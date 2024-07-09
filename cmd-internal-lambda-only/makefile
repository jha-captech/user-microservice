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


# Lambda

.PHONY: lambda_build
lambda_build:
	sam build

.PHONY: lambda_local_api
lambda_local_api: db_up lambda_build
	sam local start-api -p 8080 --env-vars env.json