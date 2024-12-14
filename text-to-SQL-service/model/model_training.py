import torch
import os
import matplotlib.pyplot as plt
from transformers import T5Tokenizer, T5ForConditionalGeneration, Trainer, TrainingArguments
from datasets import load_dataset
from sklearn.metrics import accuracy_score

class TextToSQLTrainer:
    def __init__(self, model_name='t5-large', max_input_length=512, max_target_length=256):
        self.device = torch.device("cuda" if torch.cuda.is_available() else "cpu")
        print(f"Using device: {self.device}")

        self.tokenizer = T5Tokenizer.from_pretrained(model_name)
        self.model = T5ForConditionalGeneration.from_pretrained(model_name).to(self.device)

        self.max_input_length = max_input_length
        self.max_target_length = max_target_length

    def preprocess_function(self, examples):
        inputs = [f"Translate Natural Language to SQL: {q}" for q in examples["question"]]
        targets = examples["sql"]

        model_inputs = self.tokenizer(
            inputs,
            max_length=self.max_input_length,
            truncation=True,
            padding="max_length"
        )

        labels = self.tokenizer(
            targets,
            max_length=self.max_target_length,
            truncation=True,
            padding="max_length"
        ).input_ids

        model_inputs["labels"] = labels
        return model_inputs

    def load_dataset(self):
        dataset = load_dataset("gretelai/synthetic_text_to_sql")

        tokenized_datasets = dataset.map(
            self.preprocess_function,
            batched=True,
            remove_columns=dataset["train"].column_names
        )

        self.train_dataset = tokenized_datasets["train"]
        self.test_dataset = tokenized_datasets["test"]

    def compute_accuracy(self, predictions, labels):
        pred_ids = predictions.argmax(dim=-1)
        pred_ids = pred_ids[labels != -100]
        labels = labels[labels != -100]

        return accuracy_score(labels.cpu().numpy(), pred_ids.cpu().numpy())

    def train(self):
        training_args = TrainingArguments(
            output_dir="./results",
            evaluation_strategy="epoch",
            learning_rate=3e-5,
            per_device_train_batch_size=8,
            per_device_eval_batch_size=8,
            num_train_epochs=20,
            weight_decay=0.01,
            save_total_limit=3,
            save_steps=500,
            logging_dir="./logs",
            logging_steps=100,
            load_best_model_at_end=True,
            metric_for_best_model="eval_loss",
            greater_is_better=False,
        )

        trainer = Trainer(
            model=self.model,
            args=training_args,
            train_dataset=self.train_dataset,
            eval_dataset=self.test_dataset,
            tokenizer=self.tokenizer,
            compute_metrics=self.compute_accuracy,
        )

        trainer.train()

        self._plot_accuracy(trainer)

    def save_model(self, output_dir='./text_to_sql_model'):
        os.makedirs(output_dir, exist_ok=True)
        self.model.save_pretrained(output_dir)
        self.tokenizer.save_pretrained(output_dir)
        print(f"Model saved to {output_dir}")

    def predict(self, input_text):
        self.model.eval()
        inputs = self.tokenizer(
            f"Translate Natural Language to SQL: {input_text}",
            return_tensors="pt"
        ).to(self.device)

        outputs = self.model.generate(
            inputs.input_ids,
            max_length=self.max_target_length,
            num_beams=4,
            early_stopping=True
        )

        return self.tokenizer.decode(outputs[0], skip_special_tokens=True)

    def _plot_accuracy(self, trainer):
        train_accuracies = trainer.state.log_history
        val_accuracies = [log for log in train_accuracies if "eval_accuracy" in log]

        train_accuracy_values = [log["accuracy"] for log in train_accuracies if "accuracy" in log]
        eval_accuracy_values = [log["eval_accuracy"] for log in val_accuracies]

        plt.figure(figsize=(10, 6))
        plt.plot(range(1, len(train_accuracy_values) + 1), train_accuracy_values, label="Train Accuracy")
        plt.plot(range(1, len(eval_accuracy_values) + 1), eval_accuracy_values, label="Validation Accuracy")
        plt.xlabel("Epochs")
        plt.ylabel("Accuracy")
        plt.title("Training and Validation Accuracy")
        plt.legend()

        # Сохраняем график в файл
        plt.savefig("text_to_sql_accuracy.png")
        plt.close()

def main():
    trainer = TextToSQLTrainer(model_name='t5-large')

    trainer.load_dataset()

    trainer.train()

    trainer.save_model()

    test_query = "Get all active users from New York"
    prediction = trainer.predict(test_query)
    print("Predicted SQL:", prediction)

if __name__ == "__main__":
    main()
