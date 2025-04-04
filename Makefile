DOCKER_COMPOSE := docker compose -f compose.local.yaml

all: down build up

.PHONY: build
build:
	${DOCKER_COMPOSE} build
	cd frontend; npm install

.PHONY: up
up:
	${DOCKER_COMPOSE} up -d
	cd frontend; npm run dev

.PHONY: down
down:
	${DOCKER_COMPOSE} down --remove-orphans

.PHONY: logs
logs:
	${DOCKER_COMPOSE} logs

.PHONY: logsf
logsf:
	${DOCKER_COMPOSE} logs -f

.PHONY: psql
psql:
	${DOCKER_COMPOSE} up --wait db
	${DOCKER_COMPOSE} exec db psql --user=postgres albatross

.PHONY: psql-query
psql-query:
	${DOCKER_COMPOSE} up --wait db
	${DOCKER_COMPOSE} exec --no-TTY db psql --user=postgres albatross

.PHONY: sqldef-dryrun
sqldef-dryrun: down
	${DOCKER_COMPOSE} build db
	${DOCKER_COMPOSE} up --wait db
	${DOCKER_COMPOSE} run --no-TTY tools psqldef --dry-run < ./backend/schema.sql

.PHONY: sqldef
sqldef: down
	${DOCKER_COMPOSE} build db
	${DOCKER_COMPOSE} up --wait db
	${DOCKER_COMPOSE} run --no-TTY tools psqldef < ./backend/schema.sql

.PHONY: asynq
asynq:
	${DOCKER_COMPOSE} up --wait task-db
	${DOCKER_COMPOSE} run tools go run github.com/hibiken/asynq/tools/asynq --uri task-db:6379 dash

.PHONY: init
init: build initdb

.PHONY: initdb
initdb:
	make psql-query < ./backend/schema.sql
	make psql-query < ./backend/fixtures/dev.sql

.PHONY: gen
gen:
	cd backend; make gen
	cd frontend; npm run openapi-typescript
