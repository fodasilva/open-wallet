run-dev:
	go run cmd/api/main.go

test:
	go test -v ./...

test-migrations:
	go test -v -tags migrations tests/e2e/migrations_test.go

compose:
	docker compose	up -d

db-migrate:
	migrate -path ./migrations -database "postgres://docker:docker@localhost:5432/docker?sslmode=disable" up
	
db-migrate-create:
	migrate create -ext sql -dir migrations -seq $(name)

db-migrate-install:
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

gen-repos:
	bash scripts/repository/gen-repos.sh

gen-docs:
	swag init -g cmd/api/main.go --parseDependency --parseInternal

gen-docs-install:
	go install github.com/swaggo/swag/cmd/swag@v1.16.4

lint:
	golangci-lint run ./...

lint-install:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

check-docs:
	@swag init -g cmd/api/main.go --parseDependency --parseInternal
	@if [ -n "$$(git status -s docs/)" ]; then \
		echo "Documentation is out of sync!"; \
		echo "Please run 'swag init -g cmd/api/main.go' locally and commit the updated docs/ folder."; \
		git diff docs/; \
		exit 1; \
	fi
	@echo "Documentation is up to date."
	
check-repos:
	@bash scripts/repository/gen-repos.sh
	@if [ -n "$$(git status -s internal/resources/)" ]; then \
		echo "Repositories are out of sync!"; \
		echo "Please run 'make gen-repos' locally and commit the updated files."; \
		git status -s internal/resources/; \
		exit 1; \
	fi
	@echo "Repositories are up to date."
