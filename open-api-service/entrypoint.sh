#!/bin/sh
if [ ! -f /app/config/local.yaml ]; then
  mkdir -p /app/config
  cat > /app/config/local.yaml <<EOL
ApiGatewayURL: "api-gateway:5001"
PORT: "5009"
GATEWAY: "api-gateway:5001"
VisualisationService: "visualisation-service-container:5007"
EOL
fi

exec "$@"
