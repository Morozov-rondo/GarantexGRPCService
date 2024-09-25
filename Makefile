build:
	go build ./cmd/main.go

test:
	go test ./... -v -cover

docker-build:
	docker build --tag exchange:dev  .

run:
	docker compose up

stop:
	docker compose down

lint:
	golangci-lint run ./...

