.PHONY: up
up:
	docker-compose up -d


.PHONY: down
down:
	docker-compose down


.PHONY: http_dev
http_dev:
	go run ./cmd/http