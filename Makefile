NAME="media-indexer"



build:
	@echo "Building..."
	@go build -o bin/$(NAME) main.go

build-linux:
	@echo "Building for Linux..."
	@GOOS=linux GOARCH=amd64 go build -o bin/$(NAME) main.go

build-windows:
	@echo "Building for Windows..."
	@GOOS=windows GOARCH=amd64 go build -o bin/$(NAME).exe main.go

build-docker:
	@echo "Building Docker image..."
	@docker build -t $(NAME) .