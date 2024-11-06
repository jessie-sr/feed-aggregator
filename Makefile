build:
	@go build -o bin/rss-agg

run: build
	@./bin/rss-agg

test:
	@go test ./... -v