from transformers import T5Tokenizer, T5ForConditionalGeneration
from config import MODEL_DIR

tokenizer = T5Tokenizer.from_pretrained(MODEL_DIR)
model = T5ForConditionalGeneration.from_pretrained(MODEL_DIR)

def generate_sql_query(text_query: str) -> str:
    inputs = tokenizer("translate English to SQL: " + text_query, return_tensors="pt", padding=True)
    outputs = model.generate(**inputs, max_new_tokens=50)
    sql_query = tokenizer.decode(outputs[0], skip_special_tokens=True)
    return sql_query
