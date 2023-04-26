env:
	docker-compose start || docker-compose up -d

migrate:
	go run migration/main.go

run:
	go run -race .

test:
	go test -count=1 -race -coverprofile=coverage.out  ./... && go tool cover -func=coverage.out

fmt:
	go fmt ./...

tidy:
	go mod tidy && go mod vendor

.PHONY: env migrate run test fmt tidy
