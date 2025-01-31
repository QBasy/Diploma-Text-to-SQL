from transformers import T5ForConditionalGeneration, T5Tokenizer
from datasets import load_dataset
from transformers import TrainingArguments, Trainer
import torch

device = "cuda" if torch.cuda.is_available() else "cpu"
print(f"Using device: {device}")

model_name = "t5-base"
model = T5ForConditionalGeneration.from_pretrained(model_name)
tokenizer = T5Tokenizer.from_pretrained(model_name)

model = model.to(device)

dataset = load_dataset("spider")

def preprocess_function(examples):
    inputs = [f"Translate English to SQL: {question}" for question in examples["question"]]
    targets = examples["query"]

    model_inputs = tokenizer(inputs, max_length=256, truncation=True, padding="max_length")

    with tokenizer.as_target_tokenizer():
        labels = tokenizer(targets, max_length=256, truncation=True, padding="max_length")

    model_inputs["labels"] = labels["input_ids"]

    return model_inputs

tokenized_dataset = dataset.map(preprocess_function, batched=True, remove_columns=dataset["train"].column_names)

training_args = TrainingArguments(
    output_dir="./text-to-sql-simple",
    evaluation_strategy="steps",
    eval_steps=500,
    save_strategy="steps",
    save_steps=500,
    learning_rate=5e-5,
    per_device_train_batch_size=4,
    per_device_eval_batch_size=4,
    gradient_accumulation_steps=4,
    fp16=True,
    num_train_epochs=10,
    weight_decay=0.01,
    save_total_limit=2,
    logging_dir="./logs",
    logging_steps=50,
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

model.save_pretrained("./text-to-sql-simple")
tokenizer.save_pretrained("./text-to-sql-simple")

print("Training complete. Model saved to ./text-to-sql-simple")
