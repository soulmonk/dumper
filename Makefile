#-include .env
SELF_DIR := $(dir $(lastword $(MAKEFILE_LIST)))
include $(SELF_DIR)/cuppa-tools/Makefile.common.mk

PROJ_BUILD_PATH := $(PROJ_BASE)/build

test:
	echo ${PROJ_BASE}

migration:
	@echo "Migration"
	@migrate -database ${POSTGRESQL_URL} -path db/migration up

gen-db: ## Generate go from sql
	@echo "Generate"
	@sqlc generate

build: ## Build server
	@echo "  >  Build server"
	go build -o $(PROJ_BUILD_PATH)/server $(PROJ_BASE)/cmd/server.go
