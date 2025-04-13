from transformers import T5ForConditionalGeneration, T5Tokenizer
import torch

device = torch.device("cuda" if torch.cuda.is_available() else "cpu")

model_name = "t5-small"
model = T5ForConditionalGeneration.from_pretrained(model_name).to(device)
tokenizer = T5Tokenizer.from_pretrained(model_name)

input_text = "Translate English to SQL: покажи все записи из таблицы users"
input_ids = tokenizer(input_text, return_tensors="pt").input_ids.to(device)

with torch.no_grad():
    outputs = model.generate(input_ids)

sql_query = tokenizer.decode(outputs[0], skip_special_tokens=True)
print(sql_query)

print("2")
tokenizer = T5Tokenizer.from_pretrained(model_name)
print("3")
model.save_pretrained("./models/codet5-large")
print("4")
tokenizer.save_pretrained("./models/codet5-large")
print("5")
