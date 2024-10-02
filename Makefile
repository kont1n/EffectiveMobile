include .env

build:
	docker compose -f docker-compose.yml build

run:
	docker compose -f docker-compose.yml up

test:
	go test -v ./...
