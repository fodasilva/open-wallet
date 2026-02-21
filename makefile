run-dev:
	go run cmd/api/main.go
test:
	go test -v ./...
compose:
	docker compose	up -d
db-migrate:
	migrate -path ./migrations -database "postgres://docker:docker@localhost:5432/docker?sslmode=disable" up
db-migrate-create:
	migrate create -ext sql -dir migrations -seq $(name)
