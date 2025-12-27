#!/bin/bash

# Ensure the database is running
# docker-compose up -d db

# Run the application with local environment variables
export DB_USER=user
export DB_PASSWORD=password
export DB_HOST=localhost
export DB_NAME=books_db
export OTEL_EXPORTER_OTLP_ENDPOINT=localhost:4317

echo "Starting api..."
./api
