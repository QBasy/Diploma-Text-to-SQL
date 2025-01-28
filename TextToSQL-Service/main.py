from fastapi import FastAPI, HTTPException
from pydantic import BaseModel
from transformers import T5ForConditionalGeneration, T5Tokenizer
import torch
import logging

logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

app = FastAPI()

model_name = "t5-small"
model = T5ForConditionalGeneration.from_pretrained(model_name)
tokenizer = T5Tokenizer.from_pretrained(model_name)

device = "cuda" if torch.cuda.is_available() else "cpu"
model = model.to(device)

class QueryRequestComplex(BaseModel):
    schema: dict
    query: str

class QueryRequestSimple(BaseModel):
    query: str


@app.on_event("startup")
async def startup():
    redis = aioredis.from_url("redis://localhost:6379")
    FastAPICache.init(RedisBackend(redis), prefix="fastapi-cache")

@app.post("/text-to-sql/complex")
@cache(expire=60)
async def text_to_sql_complex(request: QueryRequestComplex):
    try:
        logger.info(f"Received request: {request.query}")

        schema_text = " ".join([f"Table {table['name']} has columns {', '.join([col['name'] for col in table['columns']])}."
                                for table in request.schema["tables"]])
        input_text = f"Schema: {schema_text}. Translate English to SQL: {request.query}"

        input_ids = tokenizer(input_text, return_tensors="pt").input_ids.to(device)
        outputs = model.generate(input_ids)
        sql_query = tokenizer.decode(outputs[0], skip_special_tokens=True)

        logger.info(f"Generated SQL: {sql_query}")
        return {"sql_query": sql_query}
    except Exception as e:
        logger.error(f"Error: {str(e)}")
        raise HTTPException(status_code=500, detail=str(e))

@app.post("/text-to-sql/simple")
async def text_to_sql_simple(request: QueryRequestSimple):
    try:
        logger.info(f"Received request: {request.query}")
        if not request.query:
            raise HTTPException(status_code=400, detail="Query cannot be empty")

        input_text = f"Translate English to SQL: {request.query}"
        input_ids = tokenizer(input_text, return_tensors="pt").input_ids.to(device)
        outputs = model.generate(input_ids)
        sql_query = tokenizer.decode(outputs[0], skip_special_tokens=True)

        if not sql_query:
            raise HTTPException(status_code=500, detail="Failed to generate SQL query")

        logger.info(f"Generated SQL: {sql_query}")
        return {"sql_query": sql_query}
    except HTTPException as e:
        raise e
    except Exception as e:
        logger.error(f"Error: {str(e)}")
        raise HTTPException(status_code=500, detail=str(e))

@app.get("/health")
async def health_check():
    return {"status": "healthy"}