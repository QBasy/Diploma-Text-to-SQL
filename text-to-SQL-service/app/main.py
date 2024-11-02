from fastapi import FastAPI, HTTPException
from models import QueryRequest, QueryResponse
from sql_generator import generate_sql_query

app = FastAPI()

@app.post("/query", response_model=QueryResponse)
async def text_to_sql(query: QueryRequest):
    try:
        sql_query = generate_sql_query(query.text)
        return QueryResponse(sql_query=sql_query)
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))