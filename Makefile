.PHONY: setup build run clean test run-scratch

setup:
	./scripts/setup-env.sh

build:
	docker-compose build

run:
	NONINTERACTIVE=true ./scripts/setup-env.sh
	docker-compose --env-file .env up -d

run-scratch:
	docker-compose down -v
	NONINTERACTIVE=true ./scripts/setup-env.sh
	docker-compose --env-file .env up -d --build

clean:
	docker-compose down -v

test:
	go test ./...