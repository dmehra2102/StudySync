.PHONY: build run test clean migrate lint

build:
		@echo "Building application..."
		go build -o bin/studysync ./cmd/api

run:
		@echo "Starting application..."
		go run ./cmd/api

test:
		@echo "Running tests..."
		go test -v ./...

clean:
		@echo "Cleaning..."
		rm -rf bin/

lint:
		@echo "Linting..."
		golangci-lint run

docker-build:
	docker-compose build

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

docker-logs:
	docker-compose logs -f
