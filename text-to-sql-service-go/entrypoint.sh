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
