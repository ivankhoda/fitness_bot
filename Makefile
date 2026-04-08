include .env

export

COMPOSE ?= docker compose
MIGRATE ?= migrate
MIGRATE_PATH ?= /migrations
MIGRATE_DB_HOST ?= db
MIGRATE_DATABASE ?= postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(MIGRATE_DB_HOST):5432/$(POSTGRES_DB)?sslmode=$(POSTGRES_SSLMODE)

.PHONY: up down logs migrate-up migrate-down migrate-force migrate-version migrate-create

up:
	$(COMPOSE) up --build

down:
	$(COMPOSE) down -v

logs:
	$(COMPOSE) logs -f api db migrate

migrate-up:
	$(COMPOSE) run --rm migrate -path $(MIGRATE_PATH) -database "$(MIGRATE_DATABASE)" up

migrate-down:
	$(COMPOSE) run --rm migrate -path $(MIGRATE_PATH) -database "$(MIGRATE_DATABASE)" down 1

migrate-force:
	@test -n "$(VERSION)" || (echo "VERSION is required" && exit 1)
	$(COMPOSE) run --rm migrate -path $(MIGRATE_PATH) -database "$(MIGRATE_DATABASE)" force $(VERSION)

migrate-version:
	$(COMPOSE) run --rm migrate -path $(MIGRATE_PATH) -database "$(MIGRATE_DATABASE)" version

migrate-create:
	@test -n "$(NAME)" || (echo "NAME is required" && exit 1)
	@command -v $(MIGRATE) >/dev/null 2>&1 || (echo "migrate CLI is required; install it with brew install golang-migrate" && exit 1)
	$(MIGRATE) create -ext sql -dir migrations -seq $(NAME)
