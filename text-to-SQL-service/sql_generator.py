# app/sql_generator.py
from transformers import AutoTokenizer, AutoModelForSeq2SeqLM

# Load a pretrained model
tokenizer = AutoTokenizer.from_pretrained("microsoft/azure-openai-codex")
model = AutoModelForSeq2SeqLM.from_pretrained("microsoft/azure-openai-codex")

def generate_sql_query(text_query: str) -> str:
    inputs = tokenizer(text_query, return_tensors="pt")
    outputs = model.generate(**inputs)
    sql_query = tokenizer.decode(outputs[0], skip_special_tokens=True)
    return sql_query
