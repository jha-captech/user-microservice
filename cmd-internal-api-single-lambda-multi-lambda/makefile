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
	docker-compose down postgres

# App

.PHONY: app_dev
app_dev:
	go run ./cmd/api

.PHONY: app
app:
	docker-compose up

.PHONY: app_up
app_up:
	docker-compose up -d

.PHONY: app_down
app_down:
	docker-compose down

# Single Lambda

.PHONY: single_lambda_build
single_lambda_build:
	sam build --template-file single_lambda.template.yaml --no-cached

.PHONY: single_lambda_local_api
single_lambda_local_api: db_up single_lambda_build
	sam local start-api -p 8080 --env-vars env.json

# Multiple Lambda

.PHONY: multiple_lambda_build
multiple_lambda_build:
	sam build --template-file multiple_lambda.template.yaml --no-cached

.PHONY: multiple_lambda_local_api
multiple_lambda_local_api: db_up multiple_lambda_build
	sam local start-api -p 8080 --env-vars env.json
	docker-compose down postgres