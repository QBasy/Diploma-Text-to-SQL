from transformers import AutoTokenizer, AutoModelForSeq2SeqLM
from config import MODEL_NAME

tokenizer = AutoTokenizer.from_pretrained(MODEL_NAME)
model = AutoModelForSeq2SeqLM.from_pretrained(MODEL_NAME)

def generate_sql_query(text_query: str) -> str:
    inputs = tokenizer(text_query, return_tensors="pt")
    outputs = model.generate(**inputs)
    sql_query = tokenizer.decode(outputs[0], skip_special_tokens=True)
    return sql_query
