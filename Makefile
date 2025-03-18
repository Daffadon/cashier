CONTROLLER_TEST_PATH =./test/ut/controller/...
SERVICE_TEST_PATH = ./test/ut/service/...
REPOSITORY_TEST_PATH = ./test/ut/repository/...
INTEGRATION_TEST_PATH?=./test/it/...

ENV_LOCAL_TEST=\
  DB_PASS="password" \
  DB_NAME="cashier" \
  DB_HOST="localhost" \
  DB_USER="user-name" \
	DB_PORT="5432" \
	APP_ENV="TEST"

all: run

build:
	@go build -o bin/main main.go

run: build
	@./bin/main

clean:
	@go clean
	@rm -rf ./bin

migrate-up:
	go run main.go migrate-up

migrate-down:
	go run main.go migrate-down

compose_up:
	@docker compose up -d --build

compose_down:
	@docker compose down

dep:
	@go mod download

vet:
	@go vet ./...

lint:
	@golangci-lint run ./... --fix

coverage_test:
	@go test -count=1 -cover -coverpkg=./controller $(CONTROLLER_TEST_PATH)
	@go test -count=1 -cover -coverpkg=./service $(SERVICE_TEST_PATH)
	@go test -count=1 -cover -coverpkg=./repository $(REPOSITORY_TEST_PATH)

unit_test:
	@go test -count=1 $(CONTROLLER_TEST_PATH)
	@go test -count=1 $(SERVICE_TEST_PATH)
	@go test -count=1 $(REPOSITORY_TEST_PATH)

integration_test:
	@$(ENV_LOCAL_TEST) \
	go test $(INTEGRATION_TEST_PATH) -count=1
