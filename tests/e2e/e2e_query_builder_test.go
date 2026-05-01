package e2e

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Masterminds/squirrel"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/felipe1496/open-wallet/infra"
	"github.com/felipe1496/open-wallet/internal/factory"
	"github.com/felipe1496/open-wallet/internal/middlewares"
	"github.com/felipe1496/open-wallet/internal/util"
	"github.com/felipe1496/open-wallet/internal/util/httputil"
	"github.com/felipe1496/open-wallet/internal/util/querybuilder"
)

type QueryBuilderTestData struct {
	ID       string  `json:"id"`
	Name     *string `json:"name"`
	Age      int     `json:"age"`
	IsActive bool    `json:"is_active"`
}

// SetupTestEngine replicates the real application's middleware stack from cmd/api/main.go
func SetupTestEngine(cfg *infra.Config, db *sql.DB, f *factory.Factory) *http.ServeMux {
	r := http.NewServeMux()
	return r
}

func RegisterQueryBuilderTestRoutes(r *http.ServeMux, resources *TestResources, cfg *infra.Config, f *factory.Factory) {
	config := querybuilder.ParseConfig{
		AllowedFields: map[string]querybuilder.FieldConfig{
			"name":      {AllowedOperators: []string{"eq", "ne", "like", "in"}},
			"age":       {AllowedOperators: []string{"eq", "ne", "gt", "gte", "lt", "lte", "in"}},
			"is_active": {AllowedOperators: []string{"eq", "in"}},
			"id":        {AllowedOperators: []string{"eq", "in"}},
		},
		AllowedSortFields: []string{"name", "age", "id"},
	}

	max, win := cfg.RateLimits.MD()

	handler := httputil.Chain(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			builder := querybuilder.Get(r.Context())

			// 1. Data Query
			query := squirrel.Select("*").From("query_builder_tests").PlaceholderFormat(squirrel.Dollar)
			query = querybuilder.ToSquirrel(query, builder)

			sqlStr, args, err := query.ToSql()
			if err != nil {
				httputil.JSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
				return
			}

			rows, err := resources.DB.QueryContext(r.Context(), sqlStr, args...)
			if err != nil {
				httputil.JSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
				return
			}
			defer func() { _ = rows.Close() }()

			var results []QueryBuilderTestData
			for rows.Next() {
				var d QueryBuilderTestData
				err := rows.Scan(&d.ID, &d.Name, &d.Age, &d.IsActive)
				if err != nil {
					httputil.JSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
					return
				}
				results = append(results, d)
			}

			// 2. Count Query (respecting filters)
			countBuilder := querybuilder.ForCount(builder)
			countQuery := squirrel.Select("COUNT(*)").From("query_builder_tests").PlaceholderFormat(squirrel.Dollar)
			countQuery = querybuilder.ToSquirrel(countQuery, countBuilder)

			countSql, countArgs, err := countQuery.ToSql()
			if err != nil {
				httputil.JSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
				return
			}

			var totalItems int
			err = resources.DB.QueryRowContext(r.Context(), countSql, countArgs...).Scan(&totalItems)
			if err != nil {
				httputil.JSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
				return
			}

			// 3. Return Paginated Response
			page := util.GetInt(r.Context(), util.ContextKeyPage)
			perPage := util.GetInt(r.Context(), util.ContextKeyPerPage)

			httputil.JSON(w, http.StatusOK, util.PaginatedResponse[[]QueryBuilderTestData]{
				Data:  results,
				Query: querybuilder.BuildMetadata(page, perPage, totalItems),
			})
		}),
		middlewares.QueryBuilderMiddleware(config),
		middlewares.NewRateLimitMiddleware(f.CacheService(), max, win, "global"),
		middlewares.CorsMiddleware(cfg),
		middlewares.DelayMiddleware(cfg),
	)

	r.Handle("GET /test-query-builder", handler)
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
		Environment: "test",
		Delay:       0,
		Origins:     []string{"http://localhost:3000"},
		RateLimits: infra.RateLimits{
			MD: func() (int, int) { return 1000, 60000 },
			XS: func() (int, int) { return 10, 60000 },
			SM: func() (int, int) { return 30, 60000 },
			LG: func() (int, int) { return 120, 60000 },
			XL: func() (int, int) { return 240, 60000 },
		},
	}

	f := factory.NewFactory(resources.DB, cfg)
	r := SetupTestEngine(cfg, resources.DB, f)
	RegisterQueryBuilderTestRoutes(r, resources, cfg, f)

	tests := []struct {
		name           string
		queryParams    string
		expectedIDs    []string
		expectedStatus int
		verifyMetadata func(*testing.T, util.PaginatedResponse[[]QueryBuilderTestData])
	}{
		{
			name:        "Metadata: Filtered Count - age eq 25",
			queryParams: "filter=age eq 25&page=1&per_page=1",
			expectedIDs: []string{"01"}, // Alice
			verifyMetadata: func(t *testing.T, res util.PaginatedResponse[[]QueryBuilderTestData]) {
				assert.Equal(t, 2, res.Query.TotalItems) // Alice and Eve
				assert.Equal(t, 2, res.Query.TotalPages)
				assert.Equal(t, true, res.Query.NextPage)
			},
		},
		{
			name:        "Metadata: Filtered Count - is_active eq true",
			queryParams: "filter=is_active eq true&per_page=10",
			expectedIDs: []string{"01", "03", "05"},
			verifyMetadata: func(t *testing.T, res util.PaginatedResponse[[]QueryBuilderTestData]) {
				assert.Equal(t, 3, res.Query.TotalItems)
				assert.Equal(t, 1, res.Query.TotalPages)
				assert.Equal(t, false, res.Query.NextPage)
			},
		},
		{
			name:        "Metadata: Last Page",
			queryParams: "filter=age eq 25&page=2&per_page=1",
			expectedIDs: []string{"05"}, // Eve
			verifyMetadata: func(t *testing.T, res util.PaginatedResponse[[]QueryBuilderTestData]) {
				assert.Equal(t, 2, res.Query.TotalItems)
				assert.Equal(t, 2, res.Query.TotalPages)
				assert.Equal(t, false, res.Query.NextPage)
			},
		},
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

			var res util.PaginatedResponse[[]QueryBuilderTestData]
			err := json.Unmarshal(w.Body.Bytes(), &res)
			require.NoError(t, err)

			var actualIDs []string
			for _, item := range res.Data {
				actualIDs = append(actualIDs, item.ID)
			}

			assert.Equal(t, tt.expectedIDs, actualIDs)

			if tt.verifyMetadata != nil {
				tt.verifyMetadata(t, res)
			}
		})
	}
}
