build:
	go build -o app cmd/app/main.go

test:
	go test ./unit_tests/...

docker-build:
	docker-compose up --build

run:
	go run cmd/app/main.go

lint:
	golangci-lint run
