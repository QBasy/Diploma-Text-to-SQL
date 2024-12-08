from fastapi import FastAPI, HTTPException
from fastapi.middleware.cors import CORSMiddleware

from pydantic import BaseModel
from transformers import AutoTokenizer, AutoModelForSeq2SeqLM
import uvicorn
import logging
from datetime import datetime
from typing import Optional, Dict

logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s',
    handlers=[logging.StreamHandler()]
)
logger = logging.getLogger(__name__)

app = FastAPI(
    title="Text-to-SQL Microservice",
    description="Convert natural language queries to SQL statements for Users and Items",
    version="1.0.0"
)

origins = [
    "http://localhost:5001",
    "*",
]

app.add_middleware(
    CORSMiddleware,
    allow_origins=origins,
    allow_credentials=True,
    allow_methods=["GET", "POST"],
    allow_headers=["*"],
)

MODEL_NAME = "cssupport/t5-small-awesome-text-to-sql"
logger.info(f"Loading model: {MODEL_NAME}")
tokenizer = AutoTokenizer.from_pretrained(MODEL_NAME)
model = AutoModelForSeq2SeqLM.from_pretrained(MODEL_NAME)
logger.info("Model and tokenizer loaded successfully.")

class QueryInput(BaseModel):
    text: str

class QueryOutput(BaseModel):
    sql_query: str
    confidence: Optional[float] = None
    execution_time: float

def convert_to_sql(text: str) -> Dict:
    try:
        start_time = datetime.now()

        input_text = f"### Question: {text}\n### SQL Query:"
        logger.info(f"Processing input: {input_text}")

        inputs = tokenizer(input_text, return_tensors="pt", max_length=512, truncation=True)
        outputs = model.generate(
            inputs.input_ids,
            max_length=128,
            num_beams=4,
            early_stopping=True
        )

        sql_query = tokenizer.decode(outputs[0], skip_special_tokens=True)

        confidence = outputs.sequences_scores[0].item() if outputs.sequences_scores else None
        execution_time = (datetime.now() - start_time).total_seconds()
        logger.info(f"Generated SQL query: {sql_query}")

        return {
            "sql_query": sql_query,
            "confidence": confidence,
            "execution_time": execution_time
        }
    except Exception as e:
        logger.error(f"Error converting text to SQL: {str(e)}")
        raise HTTPException(status_code=500, detail="Error converting text to SQL.")

@app.post("/convert", response_model=QueryOutput)
async def convert_endpoint(query_input: QueryInput):
    logger.info(f"Received query: {query_input.text}")
    result = convert_to_sql(query_input.text)
    return result

@app.get("/health")
async def health_check():
    logger.info("Health check requested.")
    return {"status": "healthy", "timestamp": datetime.now().isoformat()}

if __name__ == "__main__":
    uvicorn.run("main:app", host="0.0.0.0", port=5003, reload=True)
