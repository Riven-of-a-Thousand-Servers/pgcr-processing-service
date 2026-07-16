POSTGRES_USERNAME=root
POSTGRES_PASSWORD=password

export

.PHONY: migrate
build-migrate:
	go build -o bin/migrate cmd/migrate/main.go

.PHONY: run-migrate
run-migrate: build-migrate
	docker compose up -d
	bin/migrate

.PHONY: up
up:
	docker compose up --detach --build

.PHONY: watch
watch:
	docker compose watch processing-service
