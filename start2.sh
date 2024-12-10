#!/bin/bash

BASE_PATH="C:/Users/QoidynBasy/Desktop/Diploma-text-to-SQL"

declare -A services=(
    ["Database Service"]="cd $BASE_PATH/database-service && go run main.go"
    ["Text-to-SQL Service"]="cd $BASE_PATH/text-to-sql-service && python main.py"
    ["API Service"]="cd $BASE_PATH/API && go run main.go"
    ["Frontend"]="cd $BASE_PATH/diploma-frontend && npm run dev"
)

echo "Starting all services..."

generate_bat_file() {
    local service_name=$1
    local command=$2
    local bat_file="./${service_name}.bat"

    echo "@echo off" > "$bat_file"
    echo "title $service_name" >> "$bat_file"
    echo "$command" >> "$bat_file"

    echo "$bat_file"
}

for service_name in "${!services[@]}"; do
    echo "Generating .bat file for $service_name..."

    bat_file=$(generate_bat_file "$service_name" "${services[$service_name]}")

    echo "Starting $service_name in a new terminal..."
    start cmd /k "$bat_file"

    echo "$service_name started in a new terminal."
    sleep 2
done

echo -e "\nAll services started. Close the terminals manually to stop."
