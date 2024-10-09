include .env
MIGRATION_PATH = "./internal/migrations"

docker.up:
	@docker-compose up --build -d

docker.down:
	@docker-compose down -v

migrate:
	@migrate create -ext sql -dir $(MIGRATION_PATH) -seq $(NAME)

migrate.up:
	@migrate -path $(MIGRATION_PATH) -database $(DATABASE_URL) up $(N)

migrate.down:
	@migrate -path $(MIGRATION_PATH) -database $(DATABASE_URL) down $(N)


migrate.force:
	@migrate -path $(MIGRATION_PATH) -database $(DATABASE_URL) force $(VERSION)

PHONY: docker.up docker.down migrate migrate.up migrate.down
