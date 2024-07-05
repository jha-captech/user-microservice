include .env

# Database

.PHONY: db_up
db_up:
	docker-compose up postgres -d

.PHONY: db_down
db_down:
	docker-compose down postgres

.PHONY: db_seed
db_seed: db_up
	@echo $(DATABASE_NAME)
	@docker cp  ./postgres_setup.sql $(DATABASE_CONTAINER_NAME):/tmp/setup.sql
	@docker exec -i $(DATABASE_CONTAINER_NAME) psql -U $(DATABASE_USER) -d $(DATABASE_NAME) -f /tmp/setup.sql

# App

.PHONY: app_dev
app_dev:
	go run ./cmd/http

.PHONY: app
app:
	docker-compose up

.PHONY: app_up
app_up:
	docker-compose up -d

.PHONY: app_down
app_down:
	docker-compose down
