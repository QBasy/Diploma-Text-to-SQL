from openai import AsyncOpenAI
from dotenv import load_dotenv
from fastapi import FastAPI, HTTPException
from pydantic import BaseModel
from typing import Optional
from transformers import T5ForConditionalGeneration, T5Tokenizer
import torch
import logging
import os
import time

from huggingface_hub import login
load_dotenv()

hf_token = os.getenv("HUGGINGFACE_API_KEY")
if hf_token:
    login(token=hf_token)
else:
    raise ValueError("HUGGINGFACE_API_KEY is not set in environment variables")

logging.basicConfig(level=logging.INFO)
logger = logging.getLogger("uvicorn.error")
logging.getLogger("uvicorn").setLevel(logging.INFO)

app = FastAPI()

device = "cuda" if torch.cuda.is_available() else "cpu"

model_simple_name = "./text-to-sql-simple"
model_simple = T5ForConditionalGeneration.from_pretrained(model_simple_name).to(device)
tokenizer_simple = T5Tokenizer.from_pretrained(model_simple_name)

"""model_complex_name = "cssupport/t5-small-awesome-text-to-sql"
model_complex = T5ForConditionalGeneration.from_pretrained(model_complex_name).to(device)
tokenizer_complex = T5Tokenizer.from_pretrained(model_complex_name)"""

model_complex_name = "./models/codet5-large"
model_complex = T5ForConditionalGeneration.from_pretrained(model_complex_name).to(device)
tokenizer_complex = T5Tokenizer.from_pretrained(model_complex_name)

client = AsyncOpenAI(api_key=os.getenv("OPENAI_API_KEY"))

class QueryRequestComplex(BaseModel):
    schema: Optional[dict] = None
    query: str

class QueryRequestSimple(BaseModel):
    query: str

print("started working")

@app.post("/text-to-sql/simple")
async def text_to_sql_simple(request: QueryRequestSimple):
    logger.info(f"Received Simple Request: {request.query}")

    if not request.query:
        raise HTTPException(status_code=400, detail="Query cannot be empty")

    input_text = f"Translate English to SQL: {request.query}"
    input_ids = tokenizer_simple(input_text, return_tensors="pt").input_ids.to(device)

    with torch.no_grad():
        outputs = model_simple.generate(input_ids)

    sql_query = tokenizer_simple.decode(outputs[0], skip_special_tokens=True).replace('"', "'")

    if not sql_query:
        raise HTTPException(status_code=500, detail="Failed to generate SQL query")

    logger.info(f"Generated SQL (Simple): {sql_query}")
    return {"sql_query": sql_query}


@app.post("/text-to-sql/complex")
async def text_to_sql_complex(request: QueryRequestComplex):
    logger.info(f"Received Complex Request: {request.dict()}")

    if not request.schema or not request.query:
        raise HTTPException(status_code=400, detail="Schema and query cannot be empty")

    schema_text = []
    for table in request.schema["schema"]:
        columns = ", ".join([f"{col['name']} {col['type']}" for col in table["columns"]])
        constraints = []
        if "primaryKey" in table:
            constraints.append(f"PRIMARY KEY ({table['primaryKey']})")
        for col in table["columns"]:
            if col["isForeignKey"]:
                constraints.append(
                    f"FOREIGN KEY ({col['name']}) REFERENCES {col['referencedTable']}({col['referencedColumn']})")

        schema_text.append(
            f"CREATE TABLE {table['name']} ({columns}{', ' + ', '.join(constraints) if constraints else ''});")

    schema_text = " ".join(schema_text)

    input_text = f"""### Schema:
    {schema_text}

    ### Query:
    {request.query}

    ### SQL:
    """
    logger.info(f"Inputting data{input_text}")
    input_ids = tokenizer_complex(input_text, return_tensors="pt", padding=True, truncation=True).input_ids.to(device)

    input_ids = tokenizer_complex(input_text, return_tensors="pt", padding=True, truncation=True,
                                  max_length=512).input_ids.to(device)  # Adjust max_length

    with torch.no_grad():
        outputs = model_complex.generate(input_ids, max_length=1024)

    sql_query = tokenizer_complex.decode(outputs[0], skip_special_tokens=True).replace('"', "'")

    logger.info(f"Generated SQL (Complex): {sql_query}")
    return {"sql_query": sql_query}


@app.post("/text-to-sql/gpt")
async def text_to_sql_gpt(request: QueryRequestComplex):
    logger.info(f"Received GPT Request: {request.query}")
    start_time = time.time()

    if not request.query:
        raise HTTPException(status_code=400, detail="Query cannot be empty")

    is_simple = request.schema is None or not request.schema

    if is_simple:
        logger.info("Processing as simple request (no schema provided)")
        prompt = f"Convert this natural language query to SQL: {request.query}"
        system_message = "You are a SQL expert. Convert natural language queries to valid SQL statements. Return only the SQL code without any explanation."
    else:
        logger.info("Processing as complex request with schema")

        logger.info(f"Schema received: {request.schema}")
        schema_text = []

        if request.schema and isinstance(request.schema, dict):
            if "schema" in request.schema:
                tables = request.schema["schema"]
            elif "tables" in request.schema:
                tables = request.schema["tables"]
            else:
                tables = []
                logger.warning("Schema provided but couldn't find 'schema' or 'tables' key")

        for table in tables:
            col_descriptions = []
            for col in table.get("columns", []):
                col_desc = f"{col['name']} ({col.get('type', 'unknown')})"
                if col.get("isForeignKey"):
                    col_desc += f" [FK to {col.get('referencedTable')}.{col.get('referencedColumn')}]"
                col_descriptions.append(col_desc)

            schema_text.append(f"Table {table['name']}: {', '.join(col_descriptions)}")

        schema_description = "\n".join(schema_text)

        prompt = f"""
        Given the following database schema:

        {schema_description}

        Convert this natural language query to a valid SQL statement:
        "{request.query}"

        Return only the SQL code without any explanation.
        """
        system_message = "You are a SQL expert. Convert natural language queries to valid SQL based on the provided schema."

    try:
        openai_start = time.time()
        logger.info(f"Sending to OpenAI API with prompt type: {'simple' if is_simple else 'complex'}")

        response = await client.responses.create(
            model="gpt-3.5-turbo",
            input=[
                {"role": "system", "content": system_message},
                {"role": "user", "content": prompt}
            ],
        )

        openai_end = time.time()
        sql_query = response['choices'][0]['message']['content'].strip()

        logger.info(f"Generated SQL (GPT-3.5): {sql_query} in {openai_end - openai_start} seconds")
        logger.info(f"Total processing time: {time.time() - start_time} seconds")
        return {"sql_query": sql_query}

    except Exception as e:
        logger.error(f"OpenAI Error: {str(e)}")
        raise HTTPException(status_code=500, detail=f"Failed to generate SQL using GPT: {str(e)}")

@app.get("/health")
async def health_check():
    try:
        logger.info("Health check received, if it's not in logs, then chatgpt is dumb")
        return {"status": "healthy"}
    except Exception as e:
        logger.error(f"Error during health check: {e}")
        raise HTTPException(status_code=500, detail="Health check failed")

"""
@app.post("/text-to-sql/complex")
async def text_to_sql_complex(request: QueryRequestComplex):
    logger.info(f"Received Complex Request: {request.query}")

    if not request.schema or not request.query:
        raise HTTPException(status_code=400, detail="Schema and query cannot be empty")

    schema_text = " ".join([f"Table {table['name']} has columns {', '.join([col['name'] for col in table['columns']])}."
                            for table in request.schema["tables"]])
    input_text = f"Schema: {schema_text}. Translate English to SQL: {request.query}"

    input_ids = tokenizer(input_text, return_tensors="pt").input_ids.to(device)
    outputs = model.generate(input_ids)
    sql_query = tokenizer.decode(outputs[0], skip_special_tokens=True).replace('"', "'")

    logger.info(f"Generated SQL (Complex): {sql_query}")
    return {"sql_query": sql_query}
"""
