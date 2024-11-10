# Build the application
build:
	@mkdir -p bin  # Ensure the bin directory exists
	@go build -o bin/rss-agg

# Run the application (depends on build)
run: build
	@./bin/rss-agg

# Run database migrations
migrate:
	cd sql/schema && goose postgres "$(MIGRATION_URL)" up

# Run all tests with verbose output
test:
	@go test ./... -v
