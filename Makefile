APP_ENV ?= development

ifneq ($(APP_ENV), production)
include env/.env.development
else
include env/.env.production
endif

export

# Create the new confirm target
confirm:
	@echo -n 'Are you sure ? [y/N]' && read ans && [ $${ans:-N} = y ]

# for namespaces in large files use / as separator
# db/migrations/new not path its name it was "run" before
run/api:
	go run ./cmd/api

psql:
	psql ${GREENLIGHT_DB_DSN}

# for namespaces in large files use / as separator
# db/migrations/new not path its name it was "migration" before
db/migrations/new:
	@echo 'Creating migration files for${name}...'
	migrate create -seq -ext=.sql -dir=./migrations ${name}

# for namespaces in large files use / as separator
# db/migrations/new not path its name it was "up" before

#Include it as prerequisite
db/migrations/up: confirm
	@echo 'Running up migrations...'
	migrate -path ./migrations -database ${GREENLIGHT_DB_DSN} up