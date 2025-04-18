


echo "Cleaning Go cache..."
go clean --cache

echo "Tidying up Go modules..."
go mod tidy

echo "Running the Go project..."
go run cmd/api/main.go
