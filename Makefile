.PHONY: run migrate

run:
	@set -e; set -a; \
	if [ -f .env ]; then . ./.env; fi; \
	set +a; \
	go run ./cmd/api

migrate:
	@set -e; set -a; \
	if [ -f .env ]; then . ./.env; fi; \
	set +a; \
	: "$${DIRECT_DATABASE_URL:?Set DIRECT_DATABASE_URL in .env}"; \
	goose -dir migrations postgres "$$DIRECT_DATABASE_URL" up
