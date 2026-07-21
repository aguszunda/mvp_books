.PHONY: test coverage run run-local docker-up docker-down

test:
	go test ./... -v

coverage:
	@echo "Running tests with coverage..."
	@go test ./... -coverprofile=coverage.out > /dev/null
	@echo "Filtering out platform, mocks, and domain (interfaces)..."
	@grep -v -E "cmd/|mocks/|internal/domain|platform" coverage.out > coverage.filtered.out
	@go tool cover -func=coverage.filtered.out
	@echo "------------------------------------------------------------------"
	@echo "Real Logic Coverage:"
	@go tool cover -func=coverage.filtered.out | grep total | awk '{print $$3}'

uncovered: coverage
	@echo "------------------------------------------------------------------"
	@echo "Functions with 0% coverage (Potential missing tests):"
	@go tool cover -func=coverage.filtered.out | grep "0.0%" || echo "🎉 All logic is covered!"

run:
	DB_USER=developer DB_PASSWORD=admin DB_HOST=localhost DB_NAME=books_db go run cmd/api/main.go

run-local:
	docker compose up -d db
	@sleep 3
	DB_USER=developer DB_PASSWORD=admin DB_HOST=localhost DB_NAME=books_db go run cmd/api/main.go

docker-up:
	docker-compose up -d --build

docker-down:
	docker-compose down
