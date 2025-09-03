GOOSE_DRIVER=postgres
GOOSE_DBSTRING := "host=${DB_HOST} user=${DB_USER} password=${DB_PASSWORD} dbname=${DB_NAME} port=${DB_PORT} sslmode=disable"
GOOSE_MIGRATION_DIR := "migrations"

.PHONY: migrate-up migrate-down
migrate-up:
	goose up

migrate-down:
	goose down
