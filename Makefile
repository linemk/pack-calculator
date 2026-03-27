.PHONY: run test build lint ci docker-build docker-run compose-up

run:
	go run ./cmd/server

test:
	go test ./... -race -count=1 -v

build:
	go build -o bin/server ./cmd/server

lint:
	golangci-lint run ./...

ci: lint test build

docker-build:
	docker build -t pack-calculator .

docker-run:
	docker run -p 8080:8080 pack-calculator

compose-up:
	docker-compose up --build
