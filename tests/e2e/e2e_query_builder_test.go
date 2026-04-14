package e2e

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Masterminds/squirrel"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/felipe1496/open-wallet/infra"
	"github.com/felipe1496/open-wallet/internal/middlewares"
	"github.com/felipe1496/open-wallet/internal/utils/querybuilder"
)

type QueryBuilderTestData struct {
	ID       string  `json:"id"`
	Name     *string `json:"name"`
	Age      int     `json:"age"`
	IsActive bool    `json:"is_active"`
}

// SetupTestEngine replicates the real application's middleware stack from cmd/api/main.go
func SetupTestEngine(cfg *infra.Config, redisClient *redis.Client) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(middlewares.DelayMiddleware(cfg))
	r.Use(middlewares.CorsMiddleware(cfg))
	max, win := cfg.RateLimits.MD()
	r.Use(middlewares.NewRateLimitMiddleware(redisClient, max, win, "global"))
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	return r
}

func RegisterQueryBuilderTestRoutes(r *gin.Engine, resources *TestResources) {
	config := querybuilder.ParseConfig{
		AllowedFields: map[string]querybuilder.FieldConfig{
			"name":      {AllowedOperators: []string{"eq", "ne", "like", "in"}},
			"age":       {AllowedOperators: []string{"eq", "ne", "gt", "gte", "lt", "lte", "in"}},
			"is_active": {AllowedOperators: []string{"eq", "in"}},
			"id":        {AllowedOperators: []string{"eq", "in"}},
		},
		AllowedSortFields: []string{"name", "age", "id"},
	}

	r.GET("/test-query-builder", middlewares.QueryBuilderMiddleware(config), func(c *gin.Context) {
		builder := c.MustGet("query_builder").(*querybuilder.Builder)

		query := squirrel.Select("*").From("query_builder_tests").PlaceholderFormat(squirrel.Dollar)
		query = querybuilder.ToSquirrel(query, builder)

		sqlStr, args, err := query.ToSql()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		rows, err := resources.DB.QueryContext(c.Request.Context(), sqlStr, args...)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer func() { _ = rows.Close() }()

		var results []QueryBuilderTestData
		for rows.Next() {
			var d QueryBuilderTestData
			err := rows.Scan(&d.ID, &d.Name, &d.Age, &d.IsActive)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			results = append(results, d)
		}

		perPage := c.GetInt("per_page")
		if len(results) > perPage {
			results = results[:perPage]
		}

		c.JSON(http.StatusOK, results)
	})
}

func stringPtr(s string) *string { return &s }

func TestQueryBuilderE2E(t *testing.T) {
	resources := SetupTestResources(t)
	db := resources.DB

	// Create test table
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS query_builder_tests (
			id TEXT PRIMARY KEY,
			name TEXT,
			age INT,
			is_active BOOLEAN
		)
	`)
	require.NoError(t, err)

	// Seed data
	seedData := []QueryBuilderTestData{
		{ID: "01", Name: stringPtr("Alice"), Age: 25, IsActive: true},
		{ID: "02", Name: stringPtr("Bob"), Age: 30, IsActive: false},
		{ID: "03", Name: stringPtr("Charlie"), Age: 35, IsActive: true},
		{ID: "04", Name: stringPtr("David"), Age: 40, IsActive: false},
		{ID: "05", Name: stringPtr("Eve"), Age: 25, IsActive: true},
	}
	// Let's add real null for ID 06
	_, err = db.Exec("INSERT INTO query_builder_tests (id, name, age, is_active) VALUES ($1, $2, $3, $4) ON CONFLICT (id) DO UPDATE SET name = EXCLUDED.name, age = EXCLUDED.age, is_active = EXCLUDED.is_active",
		"06", nil, 0, false)
	require.NoError(t, err)

	for _, d := range seedData {
		_, err := db.Exec("INSERT INTO query_builder_tests (id, name, age, is_active) VALUES ($1, $2, $3, $4) ON CONFLICT (id) DO UPDATE SET name = EXCLUDED.name, age = EXCLUDED.age, is_active = EXCLUDED.is_active",
			d.ID, d.Name, d.Age, d.IsActive)
		require.NoError(t, err)
	}

	cfg := &infra.Config{
		Environment:    "test",
		Delay:          0,
		Origins:        []string{"http://localhost:3000"},
		RateLimitDBURL: resources.RedisConnStr,
		RateLimits: infra.RateLimits{
			MD: func() (int, int) { return 1000, 60000 },
			XS: func() (int, int) { return 10, 60000 },
			SM: func() (int, int) { return 30, 60000 },
			LG: func() (int, int) { return 120, 60000 },
			XL: func() (int, int) { return 240, 60000 },
		},
	}

	r := SetupTestEngine(cfg, resources.RedisClient)
	RegisterQueryBuilderTestRoutes(r, resources)

	tests := []struct {
		name           string
		queryParams    string
		expectedIDs    []string
		expectedStatus int
	}{
		{
			name:        "Query by name eq",
			queryParams: "filter=name eq 'Alice'",
			expectedIDs: []string{"01"},
		},
		{
			name:        "Query by age gt",
			queryParams: "filter=age gt 30",
			expectedIDs: []string{"03", "04"},
		},
		{
			name:        "Query by is_active eq true",
			queryParams: "filter=is_active eq true",
			expectedIDs: []string{"01", "03", "05"},
		},
		{
			name:        "Query LIKE case insensitive",
			queryParams: "filter=name like 'ali'",
			expectedIDs: []string{"01"},
		},
		{
			name:        "AND condition",
			queryParams: "filter=age eq 25 and is_active eq true",
			expectedIDs: []string{"01", "05"},
		},
		{
			name:        "OR group",
			queryParams: "filter=(name eq 'Alice' or name eq 'Bob')",
			expectedIDs: []string{"01", "02"},
		},
		{
			name:        "Order by age desc",
			queryParams: "order_by=age:desc,id:asc",
			expectedIDs: []string{"04", "03", "02", "01", "05", "06"},
		},
		{
			name:        "Pagination page 2 per_page 2",
			queryParams: "page=2&per_page=2&order_by=id:asc",
			expectedIDs: []string{"03", "04"},
		},
		{
			name:        "Query with IN operator (mixed types)",
			queryParams: "filter=age in (25, 35, 40)",
			expectedIDs: []string{"01", "03", "04", "05"},
		},
		{
			name:        "Query with IN operator (strings)",
			queryParams: "filter=name in ('Alice', 'Bob')",
			expectedIDs: []string{"01", "02"},
		},
		{
			name:        "Query with IN operator (booleans)",
			queryParams: "filter=is_active in (true)",
			expectedIDs: []string{"01", "03", "05"},
		},
		{
			name:        "Query with IN operator (null)",
			queryParams: "filter=name in ('Alice', null)&order_by=id:asc",
			expectedIDs: []string{"01", "06"},
		},
		{
			name:        "Query with IN operator (only null)",
			queryParams: "filter=name in (null)&order_by=id:asc",
			expectedIDs: []string{"06"},
		},
		{
			name:           "Disallowed field",
			queryParams:    "filter=unsupported_field eq '01'",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Disallowed operator",
			queryParams:    "filter=is_active gt true",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Disallowed operator inside OR group",
			queryParams:    "filter=(name eq 'Alice' or is_active gt true)",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Malformed: missing value",
			queryParams:    "filter=name eq",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Malformed: unclosed parenthesis",
			queryParams:    "filter=(name eq 'Alice'",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Malformed: multiple ANDs",
			queryParams:    "filter=name eq 'Alice' and and age gt 20",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Malformed: invalid operator",
			queryParams:    "filter=name not_an_op 'Alice'",
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/test-query-builder", nil)
			if tt.queryParams != "" {
				req.URL.RawQuery = tt.queryParams
				req.RequestURI = "/test-query-builder?" + tt.queryParams
			}

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if tt.expectedStatus != 0 {
				assert.Equal(t, tt.expectedStatus, w.Code)
				return
			}

			assert.Equal(t, http.StatusOK, w.Code)

			var results []QueryBuilderTestData
			err := json.Unmarshal(w.Body.Bytes(), &results)
			require.NoError(t, err)

			var actualIDs []string
			for _, r := range results {
				actualIDs = append(actualIDs, r.ID)
			}

			assert.Equal(t, tt.expectedIDs, actualIDs)
		})
	}
}
