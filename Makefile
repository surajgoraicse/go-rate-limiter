
MAIN="cmd/main.go"
run:
	@echo "running the app..."
	@go run $(MAIN)

test:
	@go test ./...