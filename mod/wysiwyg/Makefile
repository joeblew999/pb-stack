# Replace with your actual PocketBase URL and admin credentials if needed by custom code
export POCKETBASE_URL="http://127.0.0.1:8090"
export POCKETBASE_ADMIN_EMAIL="your_admin_email@example.com"
export POCKETBASE_ADMIN_PASSWORD="your_admin_password"

# Default flags for PocketBase serve command
PB_FLAGS?=--http="127.0.0.1:8090" --dir="./pb_data" --publicDir="./pb_public"

print:
	@echo "POCKETBASE_URL=${POCKETBASE_URL}"
	@echo "POCKETBASE_ADMIN_EMAIL=${POCKETBASE_ADMIN_EMAIL}"
	@echo "POCKETBASE_ADMIN_PASSWORD=[set or not set]" # Avoid printing password
	@echo "---"
	@echo "To run:"
	@echo "  make run"
	@echo "  make migrate-up"
	@echo "  make migrate-down"
	@echo "  make create-admin"

run: ## Start the PocketBase server with custom code
	@echo "Starting PocketBase with flags: ${PB_FLAGS}"
	# Pass GORUN_WD so main.go can potentially find the right data dir path
	GORUN_WD=$(shell pwd) go run . serve ${PB_FLAGS}

migrate-up: ## Apply pending migrations
	GORUN_WD=$(shell pwd) go run . migrate up ${PB_FLAGS}

migrate-down: ## Rollback the last migration group
	GORUN_WD=$(shell pwd) go run . migrate down ${PB_FLAGS}

create-admin: ## Create initial admin (interactive)
	GORUN_WD=$(shell pwd) go run . admin create --email="${POCKETBASE_ADMIN_EMAIL}" --password="${POCKETBASE_ADMIN_PASSWORD}" ${PB_FLAGS}

.PHONY: print run migrate-up migrate-down create-admin

