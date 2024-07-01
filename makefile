include .env


.PHONY: http_dev
http_dev:
	go run ./cmd/http


.PHONY: up
up:
	docker-compose up -d


.PHONY: down
down:
	docker-compose down


.PHONY: seed_db
seed_db: up
	@echo $(DATABASE_NAME)
	@docker cp  ./.infastructure/postgres_setup.sql $(DATABASE_CONTAINER_NAME):/tmp/setup.sql
	@docker exec -i $(DATABASE_CONTAINER_NAME) psql -U $(DATABASE_USER) -d $(DATABASE_NAME) -f /tmp/setup.sql
