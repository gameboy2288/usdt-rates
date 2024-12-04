build:
	go build -o app cmd/app/main.go

test:
	go test ./...

docker-build:
	docker build -t usdt-rates .

run:
	go run cmd/app/main.go

lint:
	golangci-lint run
