import openai
from fastapi import FastAPI, HTTPException
from pydantic import BaseModel
from transformers import T5ForConditionalGeneration, T5Tokenizer
import torch
import logging
import os
import time

logging.basicConfig(level=logging.INFO)
logger = logging.getLogger("uvicorn.error")
logging.getLogger("uvicorn").setLevel(logging.INFO)

app = FastAPI()

device = "cuda" if torch.cuda.is_available() else "cpu"

model_simple_name = "./text-to-sql-simple"
model_simple = T5ForConditionalGeneration.from_pretrained(model_simple_name).to(device)
tokenizer_simple = T5Tokenizer.from_pretrained(model_simple_name)

model_complex_name = "cssupport/t5-small-awesome-text-to-sql"
model_complex = T5ForConditionalGeneration.from_pretrained(model_complex_name).to(device)
tokenizer_complex = T5Tokenizer.from_pretrained(model_complex_name)

openai.api_key = os.getenv("OPENAI_API_KEY")

class QueryRequestComplex(BaseModel):
    schema: dict
    query: str

class QueryRequestSimple(BaseModel):
    query: str

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
    logger.info(f"Received Complex Request: {request.query}")

    if not request.schema or not request.query:
        raise HTTPException(status_code=400, detail="Schema and query cannot be empty")

    schema_text = []
    for table in request.schema["tables"]:
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

    input_text = f"""
    tables:
    {schema_text}
    query for: {request.query}
    """

    input_ids = tokenizer_complex(input_text, return_tensors="pt", padding=True, truncation=True).input_ids.to(device)

    with torch.no_grad():
        outputs = model_complex.generate(input_ids, max_length=512)

    sql_query = tokenizer_complex.decode(outputs[0], skip_special_tokens=True).replace('"', "'")

    logger.info(f"Generated SQL (Complex): {sql_query}")
    return {"sql_query": sql_query}


@app.post("/text-to-sql/gpt")
async def text_to_sql_gpt(request: QueryRequestComplex):
    logger.info(f"Received GPT Request: {request.query}")
    start_time = time.time()

    if not request.schema or not request.query:
        raise HTTPException(status_code=400, detail="Schema and query cannot be empty")

    schema_text = " ".join([f"Table {table['name']} has columns {', '.join([col['name'] for col in table['columns']])}."
                            for table in request.schema["tables"]])
    prompt = f"Schema: {schema_text}. Translate English to SQL: {request.query}"

    try:
        openai_start = time.time()
        logger.info(f"Sending to OpenAI API with prompt: {prompt}")

        response = openai.completions.create(
            model="gpt-3.5-turbo",
            prompt=prompt,
            max_tokens=2048,
            temperature=1,
            top_p=1,
            frequency_penalty=0,
            presence_penalty=0
        )

        openai_end = time.time()
        sql_query = response['choices'][0]['text'].strip()

        logger.info(f"Generated SQL (GPT-3.5): {sql_query} in {openai_end - openai_start} seconds")
        logger.info(f"Total processing time: {time.time() - start_time} seconds")
        return {"sql_query": sql_query}

    except Exception as e:
        logger.error(f"OpenAI Error: {str(e)}")
        raise HTTPException(status_code=500, detail="Failed to generate SQL using GPT")

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