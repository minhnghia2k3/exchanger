include .env
MIGRATION_PATH = "./cmd/migrate/migrations"

docker.up:
	@docker-compose up -d
	@docker logs -f exchanger-go

docker.down:
	@docker-compose down

migrate:
	@migrate create -ext sql -dir $(MIGRATION_PATH) -seq $(NAME)

migrate.up:
	@migrate -path $(MIGRATION_PATH) -database $(DATABASE_URL) up $(N)

migrate.down:
	@migrate -path $(MIGRATION_PATH) -database $(DATABASE_URL) down $(N)


migrate.force:
	@migrate -path $(MIGRATION_PATH) -database $(DATABASE_URL) force $(VERSION)

swag:
	@swag fmt --exclude docs,scripts && swag init -d ./cmd/api --parseDependency

seed:
	@go run ./cmd/migrate/seed/


PHONY: docker.up docker.down migrate migrate.up migrate.down swag seed