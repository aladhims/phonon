.PHONY: all build up down logs test clean rebuild setup

all: build up

build:
	docker-compose build

setup:
	./scripts/setup-env.sh

up: setup
	docker-compose --env-file .env up -d

down:
	docker-compose down

logs:
	docker-compose logs -f

test:
	go test ./...

clean:
	docker-compose down -v
	rebuild: down clean build up

# Service-specific commands
.PHONY: build-phonon build-cleanup restart-phonon restart-cleanup logs-phonon logs-cleanup

build-phonon:
	docker-compose build phonon

build-cleanup:
	docker-compose build cleanup

restart-phonon: build-phonon
	docker-compose restart phonon

restart-cleanup: build-cleanup
	docker-compose restart cleanup

logs-phonon:
	docker-compose logs -f phonon

logs-cleanup:
	docker-compose logs -f cleanup

# Development helpers
.PHONY: dev-deps dev-clean

dev-deps:
	go mod download
	go mod tidy

dev-clean:
	go clean
	rm -f phonon cleanup