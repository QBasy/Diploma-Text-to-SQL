# Используем официальный образ Python
FROM python:3.12-slim

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем файл с зависимостями в контейнер
COPY requirements.txt .

# Устанавливаем зависимости из requirements.txt
RUN pip install --no-cache-dir -r requirements.txt

# Копируем весь исходный код в рабочую директорию
COPY . .

# Открываем порт для FastAPI
EXPOSE 5003

# Запускаем приложение FastAPI
CMD ["uvicorn", "main:app", "--host", "0.0.0.0", "--port", "5003"]
