.PHONY: migrate-up migrate-down migrate-status migrate-create

migrate-up:
	go run cmd/migrate/main.go up

migrate-down:
	go run cmd/migrate/main.go down

migrate-status:
	go run cmd/migrate/main.go status

# Example: make migrate-create name=create_users_table
migrate-create:
	go run cmd/migrate/main.go create $(name) sql
