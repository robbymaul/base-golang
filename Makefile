############
# DIRECTORY#
############
PROJECT_ROOT ?=$(shell pwd)
WORK_DIR ?=$(PROJECT_ROOT)
PROJECT_SLUG :=payment-service

MIGRATE_DIR :=$(WORK_DIR)/migrations
PROJECT_CONFIG ?=.env

RELEASE_OUT_DIR :=$(WORK_DIR)/bin/release
RELEASE_BIN :=$(RELEASE_OUT_DIR)/$(PROJECT_SLUG)

-include $(PROJECT_ROOT)/$(PROJECT_CONFIG)


##########
# COMMON #
##########

## HELPER

## help: Show command help
.PHONY: help
all: help
help: Makefile
	@echo "Usage: make [target]"
	@echo "Choose a command run in "${PROJECT_NAME}":"
	@echo
	@sed -n '/^## /{s///;p}' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'
	@echo

#############
# CONFIGURE #
#############
######
# GO #
######
## configure: configure application environment needed
.PHONY: configure
configure: download-golang-migrate go-mod-tidy
	@go get -u gorm.io/gorm
	@go get -u gorm.io/driver/postgres
	@go get -u github.com/gin-gonic/gin
	@go get -u github.com/rs/zerolog/log

.PHONY: go-mod-tidy
go-mod-tidy:
	@echo " > go mod tidy ...."
	@go mod tidy
	@echo " > go mod tidy done ...."

## vendor : download dependencies to vendor folder
vendor: go-mod-tidy go.mod
	@-echo " > vendor : installing dependencies"
	@go mod vendor
	@-echo " > vendor : done"

###########
# RELEASE #
###########
## #
## release-bin: release application binary
.PHONY: release-bin
release-bin: vendor
	@-echo " > compiling for release"
	@CGO_ENABLED=0 go build -a -v -mod=vendor \
		-o ${RELEASE_BIN} ${WORK_DIR}/cmd/${PROJECT_SLUG}
	@-echo " > release done"
	@-${RELEASE_BIN} --version

##########
# DOCKER #
##########

## docker-local-up: running docker compose local
.PHONY: docker-local-up
docker-local-up:
	@echo " > preparation docker ..."
	@docker compose -f ${DOCKER_COMPOSE_NAME} up --build --force-recreate


## docker-local-down: stop docker compose local
.PHONY: docker-local-down
docker-local-down:
	@echo " > docker down"
	@docker compose -f ${DOCKER_COMPOSE_NAME} down

############
# DATABASE #
############
## #

MIGRATION_URL:="postgres://${MIGRATION_DB_USER}:${MIGRATION_DB_PASSWORD}@${MIGRATION_DB_HOST}:${MIGRATION_DB_PORT}/${MIGRATION_DB_NAME}?sslmode=${MIGRATION_DB_SSL_MODE}"
MIGRATION_EXECUTE:= migrate -path ${MIGRATE_DIR} -database ${MIGRATION_URL}  -verbose
## download-golang-migrate: download golang-migrate binary for database migration
.PHONY: download-golang-migrate
download-golang-migrate:
	@-echo "download golang-migrate ...."
	@go install -tags "postgres" github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	@-echo "download golang-migrate done ...."

## migrate-script: create file sql migrate script
.PHONY: migrate-script
migrate-script:
	@read -p " > migrate-script: enter migration name: " MIGRATE_NAME; \
	migrate create -ext sql -dir ${MIGRATE_DIR} -seq $$MIGRATE_NAME

## migrate-up: migrate up
.PHONY: migrate-up
migrate-up:
	@echo " > migrate up ...."
	@${MIGRATION_EXECUTE} up 1
	@echo " > migrate up done ...."

## migrate-up: migrate up all
.PHONY: migrate-up-all
migrate-up-all:
	@echo " > migrate up ...."
	@${MIGRATION_EXECUTE} up
	@echo " > migrate up done ...."

## migrate-down: migrate database down 1 step
.PHONY: migrate-down
migrate-down:
	@echo " > migrate down ...."
	@${MIGRATION_EXECUTE} down 1
	@echo " > migrate down done ...."

## migrate-clean: clean database
.PHONY: migrate-clean
migrate-clean:
	@echo " > migrate clean ...."
	@${MIGRATION_EXECUTE} drop
	@echo " > migrate clean done ...."

## migrate-version: migrate version
.PHONY: migrate-version
migrate-version:
	@echo " > migrate version ...."
	@${MIGRATION_EXECUTE} version
	@echo " > migrate version done ...."


## create-admin-auth: create user super admin
create-admin-auth: release-bin
	@echo " > db:insert-user-admin ..."
	@${RELEASE_BIN} -create-user-admin=true

## running mas dion k-wallet migration
.PHONY: mas-dion-k-wallet-migration
k-wallet-migration:
	@echo " > mas dion k-wallet migration ..."
	@${RELEASE_BIN} -migrasi-kwallet-mas-dion=true

## running mas ammar k-wallet migration
.PHONY: mas-ammar-k-wallet-migration
k-wallet-migration:
	@echo " > mas ammar k-wallet migration ..."
	@${RELEASE_BIN} -migrasi-kwallet-mas-ammar=true

## running k6 health check
.PHONY: k6-health-check
k6-health-check: 
	@read -p " > options: enter option flag: " OPTION_FLAG; 
	@echo " > k6 health check ..."
	@k6 run $$OPTION_FLAG k6/health.js

## running k6 payment scenario
.PHONY: k6-payment-scenario
k6-payment-scenario: 
	@read -p " > options: enter option flag: " OPTION_FLAG; 
	@echo " > k6 payment scenario ..."
	@k6 run $$OPTION_FLAG k6/payment/payment.js

## running k6 payment scenario
.PHONY: k6-payment-scenario-dev
k6-payment-scenario-dev:
	@read -p " > options: enter option flag: " OPTION_FLAG;
	@echo " > k6 payment scenario ..."
	@k6 run $$OPTION_FLAG k6/payment/payment_dev.js

## running k6 payment scenario
.PHONY: k6-payment-scenario-prod
k6-payment-scenario-prod:
	@read -p " > options: enter option flag: " OPTION_FLAG;
	@echo " > k6 payment scenario ..."
	@k6 run $$OPTION_FLAG k6/payment/payment_prod.js
