NAME="media-indexer"



build:
	@echo "Building..."
	@go build -o bin/$(NAME) main.go