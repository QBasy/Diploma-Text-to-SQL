from pydantic import BaseModel

class QueryRequest(BaseModel):
    text: str

class QueryResponse(BaseModel):
    sql_query: str
