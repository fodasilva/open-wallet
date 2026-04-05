package e2e

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/felipe1496/open-wallet/infra"
	"github.com/felipe1496/open-wallet/internal/services"

	"github.com/docker/go-connections/nat"
	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/oklog/ulid/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type TestResources struct {
	PostgresContainer testcontainers.Container
	RedisContainer    testcontainers.Container
	DB                *sql.DB
	RedisClient       *redis.Client
	PostgresConnStr   string
	RedisConnStr      string
}

func SetupTestResources(t *testing.T) *TestResources {
	ctx := context.Background()

	// 1. Setup Postgres
	postgresReq := testcontainers.ContainerRequest{
		Image:        "postgres:16",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     "test",
			"POSTGRES_PASSWORD": "test",
			"POSTGRES_DB":       "testdb",
		},
		WaitingFor: wait.ForListeningPort("5432/tcp"),
	}
	pgContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: postgresReq,
		Started:          true,
	})

	if err != nil {
		t.Fatalf("failed to start postgres container: %v", err)
	}

	pgHost, _ := pgContainer.Host(ctx)
	pgPort, _ := pgContainer.MappedPort(ctx, nat.Port("5432/tcp"))
	pgConnStr := fmt.Sprintf("postgres://test:test@%s:%s/testdb?sslmode=disable", pgHost, pgPort.Port())

	dbConn, err := infra.DBConn(pgConnStr)
	if err != nil {
		t.Fatalf("cannot init DB using infra.DBConn: %v", err)
	}

	if err := runMigrations(t, dbConn); err != nil {
		t.Fatalf("failed to run migrations: %v", err)
	}

	// 2. Setup Redis
	redisReq := testcontainers.ContainerRequest{
		Image:        "redis:7-alpine",
		ExposedPorts: []string{"6379/tcp"},
		WaitingFor:   wait.ForListeningPort("6379/tcp"),
	}
	redisContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: redisReq,
		Started:          true,
	})

	if err != nil {
		t.Fatalf("failed to start redis container: %v", err)
	}

	redisHost, _ := redisContainer.Host(ctx)
	redisPort, _ := redisContainer.MappedPort(ctx, nat.Port("6379/tcp"))
	redisConnStr := fmt.Sprintf("redis://%s:%s", redisHost, redisPort.Port())

	redisClient, err := infra.RedisConn(redisConnStr)
	if err != nil {
		t.Fatalf("failed to init Redis using infra.RedisConn: %v", err)
	}

	return &TestResources{
		PostgresContainer: pgContainer,
		RedisContainer:    redisContainer,
		DB:                dbConn,
		RedisClient:       redisClient,
		PostgresConnStr:   pgConnStr,
		RedisConnStr:      redisConnStr,
	}
}

// Deprecated: Use SetupTestResources instead
func SetupTestDB(t *testing.T) *PostgresTestDB {
	resources := SetupTestResources(t)
	return &PostgresTestDB{
		Container:  resources.PostgresContainer,
		DB:         resources.DB,
		ConnString: resources.PostgresConnStr,
	}
}

type PostgresTestDB struct {
	Container  testcontainers.Container
	DB         *sql.DB
	ConnString string
}

func runMigrations(_ *testing.T, db *sql.DB) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("failed to create postgres driver: %w", err)
	}

	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get working directory: %w", err)
	}

	projectRoot := wd
	for {
		if _, err := os.Stat(filepath.Join(projectRoot, "go.mod")); err == nil {
			break
		}
		parent := filepath.Dir(projectRoot)
		if parent == projectRoot {
			return fmt.Errorf("go.mod not found in directory tree")
		}
		projectRoot = parent
	}

	migrationsPath := filepath.Join(projectRoot, "migrations")
	absPath, err := filepath.Abs(migrationsPath)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", absPath),
		"postgres",
		driver,
	)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}

// TestUser represents a user created for E2E tests.
type TestUser struct {
	ID       string
	Name     string
	Email    string
	Username string
}

// SetupTestUser creates a test user in the database and returns the User object and a valid JWT token.
func SetupTestUser(t *testing.T, db *sql.DB, cfg *infra.Config) (TestUser, string) {
	user := TestUser{
		ID:       ulid.Make().String(),
		Name:     "Test User",
		Email:    "test@example.com",
		Username: "testuser",
	}

	// Create user in DB to satisfy foreign keys
	_, err := db.Exec("INSERT INTO users (id, name, email, username) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING",
		user.ID, user.Name, user.Email, user.Username)
	assert.NoError(t, err)

	jwtService := services.NewJWTService(cfg)
	token, err := jwtService.GenerateToken(user.ID)
	assert.NoError(t, err)

	return user, token
}

// AssertTableIsEmpty verifies that a given table is empty.
func AssertTableIsEmpty(t *testing.T, db *sql.DB, tableName string) {
	var count int
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s", tableName)
	err := db.QueryRow(query).Scan(&count)
	assert.NoError(t, err)
	assert.Equal(t, 0, count, "Table %s should be empty at the start of the test", tableName)
}

// AssertUnauthorized verifies that a request to a protected endpoint returns 401 when unauthenticated.
func AssertUnauthorized(t *testing.T, engine *gin.Engine, method string, url string, body io.Reader) {
	req := httptest.NewRequest(method, url, body)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code, "Endpoint %s %s should require authentication", method, url)
}
