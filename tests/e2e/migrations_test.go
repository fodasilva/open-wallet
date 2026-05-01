//go:build migrations

package e2e

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/docker/go-connections/nat"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/felipe1496/open-wallet/infra"
)

func TestMigrationsIntegrity(t *testing.T) {
	ctx := context.Background()

	// 1. Setup a fresh Postgres container
	postgresReq := testcontainers.ContainerRequest{
		Image:        "postgres:16",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     "test",
			"POSTGRES_PASSWORD": "test",
			"POSTGRES_DB":       "migration_test_db",
		},
		WaitingFor: wait.ForListeningPort("5432/tcp"),
	}
	pgContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: postgresReq,
		Started:          true,
	})
	assert.NoError(t, err)
	defer func() { _ = pgContainer.Terminate(ctx) }()

	pgHost, _ := pgContainer.Host(ctx)
	pgPort, _ := pgContainer.MappedPort(ctx, nat.Port("5432/tcp"))
	pgConnStr := fmt.Sprintf("postgres://test:test@%s:%s/migration_test_db?sslmode=disable", pgHost, pgPort.Port())

	db, err := infra.DBConn(pgConnStr)
	assert.NoError(t, err)
	defer db.Close()

	// 2. Prepare migrate instance
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	assert.NoError(t, err)

	wd, _ := os.Getwd()
	projectRoot := wd
	for {
		if _, err := os.Stat(filepath.Join(projectRoot, "go.mod")); err == nil {
			break
		}
		projectRoot = filepath.Dir(projectRoot)
	}

	migrationsPath := fmt.Sprintf("file://%s", filepath.Join(projectRoot, "migrations"))
	m, err := migrate.NewWithDatabaseInstance(migrationsPath, "postgres", driver)
	assert.NoError(t, err)

	// 3. Test UP
	t.Log("Testing migrations UP...")
	err = m.Up()
	assert.NoError(t, err, "Migrations UP failed")

	// 4. Test DOWN
	t.Log("Testing migrations DOWN...")
	err = m.Down()
	assert.NoError(t, err, "Migrations DOWN failed (check your down.sql files)")

	// 5. Test UP again (ensure rebuild works)
	t.Log("Testing migrations UP again...")
	err = m.Up()
	assert.NoError(t, err, "Migrations UP failed on second attempt")

	t.Log("All migrations are healthy!")
}
