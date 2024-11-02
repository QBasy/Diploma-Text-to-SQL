from datasets import load_dataset
from transformers import T5Tokenizer, T5ForConditionalGeneration, Trainer, TrainingArguments

# Load the dataset
ds = load_dataset("Clinton/Text-to-sql-v1")
dataset = ds['train'].rename_column("response", "target")

# Split the dataset into train and test sets
train_test_split = dataset.train_test_split(test_size=0.1)
train_dataset = train_test_split['train']
eval_dataset = train_test_split['test']

# Initialize the tokenizer and model
tokenizer = T5Tokenizer.from_pretrained("t5-small")
model = T5ForConditionalGeneration.from_pretrained("t5-small")

# Preprocessing function
def preprocess_function(examples):
    model_inputs = tokenizer(examples['input'], max_length=128, truncation=True, padding="max_length")
    labels = tokenizer(examples['target'], max_length=128, truncation=True, padding="max_length")
    model_inputs['labels'] = labels['input_ids']
    return model_inputs

# Tokenize the datasets
tokenized_train_dataset = train_dataset.map(preprocess_function, batched=True)
tokenized_eval_dataset = eval_dataset.map(preprocess_function, batched=True)

# Set training arguments
training_args = TrainingArguments(
    output_dir="./text_to_sql_model",
    eval_strategy="steps",
    eval_steps=500,
    learning_rate=3e-5,
    per_device_train_batch_size=8,  # Adjust batch size according to your hardware
    per_device_eval_batch_size=8,
    num_train_epochs=3,
    weight_decay=0.01,
    logging_dir='./logs',
    save_total_limit=2,
)

# Create a Trainer instance
trainer = Trainer(
    model=model,
    args=training_args,
    train_dataset=tokenized_train_dataset,
    eval_dataset=tokenized_eval_dataset,
)

# Train the model
trainer.train()

# Save the model and tokenizer
model.save_pretrained("./text_to_sql_model")
tokenizer.save_pretrained("./text_to_sql_model")

print("Model trained and saved successfully!")
