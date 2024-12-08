#!/bin/bash

declare -A services=(
    ["Database Service"]="cd './database-service' && go run main.go"
    ["Text-to-SQL Service"]="python3 './text-to-sql-service/main.py'"
    ["API Service"]="cd './API' && go run main.go"
    ["Frontend"]="cd './diploma-frontend' && npm run dev"
)

pids=()

stop_services() {
    echo -e "\nStopping all services..."
    for pid in "${pids[@]}"; do
        echo "Stopping process with PID $pid..."
        kill "$pid" 2>/dev/null || echo "Failed to stop process $pid."
    done
    echo "All services stopped."
    exit 0
}

trap stop_services SIGINT

echo "Starting all services..."

for service_name in "${!services[@]}"; do
    echo "Starting $service_name..."

    eval "${services[$service_name]}" > "./$service_name.log" 2>&1 &

    pid=$!
    pids+=($pid)

    echo "$service_name started with PID $pid. Logs are being written to ./$service_name.log."

    sleep 2
done

echo -e "\nAll services started. Press Ctrl+C to stop."

wait
