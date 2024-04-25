-include .env
export

CURRENT_DIR=$(shell pwd)
APP=api-gateway
CMD_DIR=./cmd

.DEFAULT_GOAL = build

# run service
.PHONY: run
run:
	go run ${CMD_DIR}/app/main.go

# proto
.PHONY: proto-gen
proto-gen:
	./scripts/gen-proto.sh

# swagger
.PHONY: swag-gen
swag-gen:
	swag init -g api/router.go -o api/docs