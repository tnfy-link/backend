# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

# Linting
GOLINT=golangci-lint

AIR=air

all: lint test benchmark

air:
	$(AIR) -c .air.toml

test:
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html

benchmark:
	$(GOTEST) -v -bench=. -benchmem ./...

lint:
	$(GOLINT) run

clean:
	$(GOCLEAN)
	rm -f coverage.out
	rm -f coverage.html

run:
	docker compose up --build

api-docs:
	swag fmt -g ./main.go \
		&& swag init --parseDependency -g ./main.go -o ./api

deps:
	$(GOGET) -v -d ./...
	go install github.com/cosmtrek/air@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/swaggo/swag/cmd/swag@latest

.PHONY: all air test benchmark lint clean run
