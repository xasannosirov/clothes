-include .env
export

CURRENT_DIR=$(shell pwd)
APP=product_service
CMD_DIR=./cmd

.PHONY: run
run:
	go run ${CMD_DIR}/app/main.go

.PHONY: proto-gen
proto-gen:
	./scripts/gen-proto.sh

.PHONY: create-migration
create-migration:
	migrate create -ext sql -dir migrations -seq "$(name)"

.PHONY: migrate-up
migrate-up:
	migrate -source file://migrations -database postgresql://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DATABASE}?sslmode=disable up

.PHONY: migrate-down
migrate-down:
	migrate -source file://migrations -database postgresql://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DATABASE}?sslmode=disable down

.PHONY: migration-version
migration-version:
	migrate -database file://migrations - database postgresql://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DATABASE}?sslmode=disable -path migrations version 

.PHONY: migrate-dirty
migrate-dirty:
	migrate -path ./migrations/ -database file://migrations - database postgresql://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DATABASE}?sslmode=disable force "$(number)"
