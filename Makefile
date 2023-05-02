NAME="media-indexer"


doc:
	@echo "==> Running godoc"
	@go install github.com/swaggo/swag/cmd/swag@latest
	swag init

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
	@docker build -t $(NAME):latest .

db-compose-up:
	@echo "==> Running docker-compose up"
	@docker compose -f docker-compose.db.yml up -d

db-compose-down:
	@echo "==> Running docker-compose down"
	@docker compose -f docker-compose.db.yml down --remove-orphans

db-compose-down-remove-volumes:
	@echo "==> Running docker-compose down --volumes"
	@docker compose -f docker-compose.db.yml down --remove-orphans --volumes