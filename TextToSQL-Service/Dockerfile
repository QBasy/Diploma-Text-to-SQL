FROM python:3.11-slim

WORKDIR /app

COPY requirements.txt /app/

RUN pip install --no-cache-dir -r requirements.txt

COPY . /app/

EXPOSE 5004

CMD ["uvicorn", "app:main", "--host", "0.0.0.0", "--port", "5004"]