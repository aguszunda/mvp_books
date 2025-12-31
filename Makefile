.PHONY: test coverage

test:
	go test ./... -v

coverage:
	@echo "Running tests with coverage..."
	@go test ./... -coverprofile=coverage.out > /dev/null
	@echo "Filtering out cmd, mocks, domain, port, and infrastructure..."
	@grep -v -E "cmd/|mocks/|internal/domain|internal/port|internal/infrastructure" coverage.out > coverage.filtered.out
	@go tool cover -func=coverage.filtered.out
	@echo "------------------------------------------------------------------"
	@echo "Real Logic Coverage:"
	@go tool cover -func=coverage.filtered.out | grep total | awk '{print $$3}'

uncovered: coverage
	@echo "------------------------------------------------------------------"
	@echo "Functions with 0% coverage (Potential missing tests):"
	@go tool cover -func=coverage.filtered.out | grep "0.0%" || echo "🎉 All logic is covered!"
