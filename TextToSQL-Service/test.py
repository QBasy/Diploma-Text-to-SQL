from fastapi import FastAPI, HTTPException
from pydantic import BaseModel
from transformers import T5ForConditionalGeneration, T5Tokenizer
import torch
app = FastAPI()

device = "cuda" if torch.cuda.is_available() else "cpu"

model_simple_name = "./drone-ai"
model_simple = T5ForConditionalGeneration.from_pretrained(model_simple_name).to(device)
tokenizer_simple = T5Tokenizer.from_pretrained(model_simple_name)

class Request(BaseModel):
    query: str

@app.post("/direction")
async def direction(request: Request):
    if not request.query:
        raise HTTPException(status_code=400, detail="Query cannot be empty")

    input_text = f"Find better direction for this: {request.query}"
    input_ids = tokenizer_simple(input_text, return_tensors="pt").input_ids.to(device)

    with torch.no_grad():
        outputs = model_simple.generate(input_ids)

    result = tokenizer_simple.decode(outputs[0], skip_special_tokens=True).replace('"', "'")

    if not result:
        raise HTTPException(status_code=500, detail="Failed to generate SQL query")

    return {"result": result}