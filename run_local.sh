#!/bin/bash

# Ensure the database is running
# docker-compose up -d db

# Run the application with local environment variables
export DB_USER=root
export DB_PASSWORD=admin2023
export DB_HOST=localhost
export DB_NAME=books_db
export OTEL_EXPORTER_OTLP_ENDPOINT=localhost:3306

echo "Starting api..."
./api
