APP_ENV ?= development

ifneq ($(APP_ENV), production)
include env/.env.development
else
include env/.env.production
endif

export

run:
	go run ./cmd/api

