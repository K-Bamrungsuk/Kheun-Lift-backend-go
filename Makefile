.PHONY: run migrate seed

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

# Mock users (password123) and lifts for local dev only. Never run against production.
seed:
	@set -e; set -a; \
	if [ -f .env ]; then . ./.env; fi; \
	set +a; \
	: "$${DIRECT_DATABASE_URL:?Set DIRECT_DATABASE_URL in .env}"; \
	goose -dir seeds -table goose_seed_version postgres "$$DIRECT_DATABASE_URL" up
