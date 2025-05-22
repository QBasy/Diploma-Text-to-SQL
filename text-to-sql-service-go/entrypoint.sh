#!/bin/sh
if [ ! -f /app/config/local.yaml ]; then
  mkdir -p /app/config
  cat > /app/config/local.yaml <<EOL
TTSQL:
  API_KEY: "OPENAI KEY"
  PORT: "5006"
  GROC:
    API_KEY: "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"
    MODEL: "meta-llama/llama-4-scout-17b-16e-instruct"
    BASE_URL: "https://api.groq.com/openai/v1/chat/completions"
EOL
fi

exec "$@"
