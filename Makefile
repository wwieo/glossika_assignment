ROOT := $(realpath .)

.PHONY: up
## up: build and run cmd service in local
up: clean format docs
	go run ./cmd/glossika/main.go -config=./config/local.json

.PHONY: clean
## clean: remove old binaries and unused file
clean:
	go clean
	go clean -testcache
	go mod tidy

.PHONY: format
## format: go format
format:
	go fmt ./...

.PHONY: docs
## docs: generate swagger docs
docs:
	swag fmt
	swag init --generalInfo ./main.go --dir ./cmd/glossika,./service/api,./service/internal/model,./service/internal/errorx

.PHONY: docker-build
## docker-build: use docker-compose to build all service
docker-build:
	docker-compose build

.PHONY: docker-up
## docker-run: use docker-compose to up all service
docker-up:
	docker-compose up