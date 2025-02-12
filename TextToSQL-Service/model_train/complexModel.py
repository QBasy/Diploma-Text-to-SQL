from transformers import T5ForConditionalGeneration, T5Tokenizer
from datasets import load_dataset
from transformers import TrainingArguments, Trainer
import torch
import json
from pathlib import Path

device = "cuda" if torch.cuda.is_available() else "cpu"
print(f"Using device: {device}")

model_name = "t5-base"
model = T5ForConditionalGeneration.from_pretrained(model_name)
tokenizer = T5Tokenizer.from_pretrained(model_name)
model = model.to(device)


def get_schema_text(db_id):
    base_path = Path("spider/database")
    db_path = base_path / db_id / "tables.json"
    try:
        with open(db_path) as f:
            tables = json.load(f)
    except FileNotFoundError:
        return ""

    schema_texts = []
    for table in tables:
        table_name = table['table_name']
        columns = [col[1] for col in table['column_names']]
        schema_text = f"Table {table_name} has columns {', '.join(columns)}."
        schema_texts.append(schema_text)

    return ' '.join(schema_texts)


def preprocess_function(examples):
    db_ids = examples["db_id"]
    questions = examples["question"]
    queries = examples["query"]

    inputs = []
    for db_id, question in zip(db_ids, questions):
        schema_text = get_schema_text(db_id)
        input_text = f"Schema: {schema_text}. Translate English to SQL: {question}"
        inputs.append(input_text)

    model_inputs = tokenizer(
        inputs,
        max_length=512,
        truncation=True,
        padding="max_length",
        return_tensors="pt",
    )

    with tokenizer.as_target_tokenizer():
        labels = tokenizer(
            queries,
            max_length=512,
            truncation=True,
            padding="max_length",
            return_tensors="pt",
        )
        model_inputs["labels"] = labels["input_ids"]

    return model_inputs


dataset = load_dataset("spider")

tokenized_dataset = dataset.map(preprocess_function, batched=True)

training_args = TrainingArguments(
    output_dir="./results",
    evaluation_strategy="steps",
    eval_steps=500,
    save_strategy="steps",
    save_steps=500,
    learning_rate=5e-5,
    per_device_train_batch_size=2,
    per_device_eval_batch_size=2,
    gradient_accumulation_steps=4,
    fp16=True,
    num_train_epochs=3,
    weight_decay=0.01,
    save_total_limit=2,
    logging_dir="./logs",
    logging_steps=10,
    load_best_model_at_end=True,
    metric_for_best_model="eval_loss",
    greater_is_better=False,
)

trainer = Trainer(
    model=model,
    args=training_args,
    train_dataset=tokenized_dataset["train"],
    eval_dataset=tokenized_dataset["validation"],
)

trainer.train()

model.save_pretrained("./fine-tuned-text-to-sql")
tokenizer.save_pretrained("./fine-tuned-text-to-sql")
