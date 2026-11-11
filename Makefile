.PHONY: run test integration check build migrate seed

run:
	go run . serve

test:
	go test -race -count=1 ./...

integration:
	@test -n "$$TEST_DATABASE_URL" || (echo 'Set TEST_DATABASE_URL to a dedicated PostgreSQL test database'; exit 1)
	go test -race -count=1 -v ./store ./integration

check:
	test -z "$$(gofmt -l $$(find . -name '*.go'))"
	go vet ./...
	go test -race -count=1 ./...

build:
	go build -trimpath -o bin/thaifood .

migrate:
	go run . migrate

seed:
	go run . seed
