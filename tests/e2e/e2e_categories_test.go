package e2e

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/felipe1496/open-wallet/infra"
	"github.com/felipe1496/open-wallet/internal/factory"
	"github.com/felipe1496/open-wallet/internal/resources/categories/handlers"
	"github.com/felipe1496/open-wallet/internal/routes"
	"github.com/oklog/ulid/v2"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func setupCategoryTestServer(pg *sql.DB, redisClient *redis.Client, cfg *infra.Config) (*gin.Engine, *factory.Factory) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	f := factory.NewFactory(pg, cfg)

	// Reuse actual routing logic
	routes.SetupCategoriesRoutes(router, f, redisClient, cfg)

	return router, f
}

func TestE2eCategories(t *testing.T) {
	res := SetupTestResources(t)
	defer func() { _ = res.PostgresContainer.Terminate(context.Background()) }()
	defer func() { _ = res.RedisContainer.Terminate(context.Background()) }()
	defer func() { _ = res.DB.Close() }()

	cfg := &infra.Config{
		JWTSecret: "test-secret",
	}

	AssertTableIsEmpty(t, res.DB, "users")
	AssertTableIsEmpty(t, res.DB, "categories")

	testUser, token := SetupTestUser(t, res.DB, cfg)

	router, _ := setupCategoryTestServer(res.DB, res.RedisClient, cfg)

	t.Run("Authentication Enforcement", func(t *testing.T) {
		endpoints := []struct {
			method string
			url    string
		}{
			{http.MethodPost, "/api/v1/categories"},
			{http.MethodGet, "/api/v1/categories"},
			{http.MethodGet, "/api/v1/categories/202401"},
			{http.MethodPatch, "/api/v1/categories/some-id"},
			{http.MethodDelete, "/api/v1/categories/some-id"},
		}

		for _, e := range endpoints {
			AssertUnauthorized(t, router, e.method, e.url, nil)
		}
	})

	t.Run("POST /categories", func(t *testing.T) {
		type testCase struct {
			name           string
			payload        handlers.CreateCategoryRequest
			expectedStatus int
			validateDB     bool
		}

		cases := []testCase{
			{
				name: "should create category with valid data",
				payload: handlers.CreateCategoryRequest{
					Name:  "Food",
					Color: "#FF5733",
				},
				expectedStatus: http.StatusCreated,
				validateDB:     true,
			},
			{
				name: "should fail when name is missing",
				payload: handlers.CreateCategoryRequest{
					Color: "#FF5733",
				},
				expectedStatus: http.StatusBadRequest,
				validateDB:     false,
			},
			{
				name: "should fail when color is missing",
				payload: handlers.CreateCategoryRequest{
					Name: "Food",
				},
				expectedStatus: http.StatusBadRequest,
				validateDB:     false,
			},
			{
				name: "should fail when name length is exactly zero",
				payload: handlers.CreateCategoryRequest{
					Name:  "",
					Color: "#FF5733",
				},
				expectedStatus: http.StatusBadRequest,
				validateDB:     false,
			},
		}

		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				body, _ := json.Marshal(tc.payload)
				req := httptest.NewRequest(http.MethodPost, "/api/v1/categories", bytes.NewBuffer(body))
				req.Header.Set("Authorization", "Bearer "+token)
				req.Header.Set("Content-Type", "application/json")
				w := httptest.NewRecorder()

				router.ServeHTTP(w, req)

				if tc.expectedStatus != w.Code {
					t.Errorf("expected status %d, got %d. Body: %s", tc.expectedStatus, w.Code, w.Body.String())
				}
				assert.Equal(t, tc.expectedStatus, w.Code)

				if tc.validateDB {
					var count int
					err := res.DB.QueryRow("SELECT COUNT(*) FROM categories WHERE name = $1 AND user_id = $2", tc.payload.Name, testUser.ID).Scan(&count)
					assert.NoError(t, err)
					assert.Equal(t, 1, count)
				}
			})
		}
	})

	t.Run("GET /categories", func(t *testing.T) {
		// Seed categories
		_, err := res.DB.Exec("INSERT INTO categories (id, user_id, name, color) VALUES ($1, $2, $3, $4), ($5, $6, $7, $8)",
			ulid.Make().String(), testUser.ID, "Transport", "#3357FF",
			ulid.Make().String(), testUser.ID, "Health", "#33FF57")
		assert.NoError(t, err)

		req := httptest.NewRequest(http.MethodGet, "/api/v1/categories", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if http.StatusOK != w.Code {
			t.Errorf("expected status 200, got %d. Body: %s", w.Code, w.Body.String())
		}
		assert.Equal(t, http.StatusOK, w.Code)

		var response handlers.ListCategoriesResponse
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, len(response.Data.Categories), 2)
	})

	t.Run("GET /categories/:period", func(t *testing.T) {
		categoryID := ulid.Make().String()
		period := "202401" // YYYYMM format as required by fn_category_amount_per_period

		// Seed category
		_, err := res.DB.Exec("INSERT INTO categories (id, user_id, name, color) VALUES ($1, $2, $3, $4)",
			categoryID, testUser.ID, "Entertainment", "#FF33F6")
		assert.NoError(t, err)

		// Seed transaction and entries
		transactionID := ulid.Make().String()
		_, err = res.DB.Exec("INSERT INTO transactions (id, user_id, name, category, category_id) VALUES ($1, $2, $3, $4, $5)",
			transactionID, testUser.ID, "Netflix", "simple_expense", categoryID)
		assert.NoError(t, err)

		_, err = res.DB.Exec("INSERT INTO entries (id, transaction_id, amount, reference_date) VALUES ($1, $2, $3, $4), ($5, $6, $7, $8)",
			ulid.Make().String(), transactionID, 15.50, "2024-01-05",
			ulid.Make().String(), transactionID, 10.00, "2024-01-20")
		assert.NoError(t, err)

		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/categories/%s", period), nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if http.StatusOK != w.Code {
			t.Errorf("expected status 200, got %d. Body: %s", w.Code, w.Body.String())
		}
		assert.Equal(t, http.StatusOK, w.Code)

		var response handlers.ListCategoryAmountPerPeriodResponse
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		found := false
		for _, cat := range response.Data.Categories {
			if cat.ID == categoryID {
				assert.Equal(t, 25.50, cat.TotalAmount)
				found = true
				break
			}
		}
		assert.True(t, found, "category not found in period list")
	})

	t.Run("PATCH /categories/:category_id", func(t *testing.T) {
		categoryID := ulid.Make().String()
		_, _ = res.DB.Exec("INSERT INTO categories (id, user_id, name, color) VALUES ($1, $2, $3, $4)",
			categoryID, testUser.ID, "Old Name", "#000000")

		_, tokenB := SetupTestUser(t, res.DB, cfg)

		type testCase struct {
			name           string
			id             string
			token          string
			payload        map[string]interface{}
			expectedStatus int
		}

		cases := []testCase{
			{
				name: "should update successfully",
				id:   categoryID,
				token: token,
				payload: map[string]interface{}{
					"name": "New Name",
				},
				expectedStatus: http.StatusOK,
			},
			{
				name: "should fail when updating non-existent category",
				id:   ulid.Make().String(),
				token: token,
				payload: map[string]interface{}{
					"name": "Whatever",
				},
				expectedStatus: http.StatusNotFound,
			},
			{
				name: "should fail when updating category of another user",
				id:   categoryID,
				token: tokenB,
				payload: map[string]interface{}{
					"name": "Hacked",
				},
				expectedStatus: http.StatusNotFound,
			},
		}

		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				body, _ := json.Marshal(tc.payload)
				req := httptest.NewRequest(http.MethodPatch, fmt.Sprintf("/api/v1/categories/%s", tc.id), bytes.NewBuffer(body))
				req.Header.Set("Authorization", "Bearer "+tc.token)
				req.Header.Set("Content-Type", "application/json")
				w := httptest.NewRecorder()

				router.ServeHTTP(w, req)

				if tc.expectedStatus == http.StatusNotFound {
					assert.Contains(t, []int{http.StatusNotFound, http.StatusForbidden}, w.Code)
				} else {
					assert.Equal(t, tc.expectedStatus, w.Code)
				}

				if tc.expectedStatus == http.StatusOK {
					var dbName string
					err := res.DB.QueryRow("SELECT name FROM categories WHERE id = $1", tc.id).Scan(&dbName)
					assert.NoError(t, err)
					assert.Equal(t, tc.payload["name"], dbName)
				}
			})
		}
	})

	t.Run("DELETE /categories/:category_id", func(t *testing.T) {
		categoryID := ulid.Make().String()
		_, _ = res.DB.Exec("INSERT INTO categories (id, user_id, name, color) VALUES ($1, $2, $3, $4)",
			categoryID, testUser.ID, "To Delete", "#FFFFFF")

		_, tokenB := SetupTestUser(t, res.DB, cfg)

		type testCase struct {
			name           string
			id             string
			token          string
			expectedStatus int
		}

		cases := []testCase{
			{
				name:           "should fail when deleting category that dont exists",
				id:             ulid.Make().String(),
				token:          token,
				expectedStatus: http.StatusNotFound,
			},
			{
				name:           "should fail when deleting category of another user",
				id:             categoryID,
				token:          tokenB,
				expectedStatus: http.StatusNotFound,
			},
			{
				name:           "should delete successfully",
				id:             categoryID,
				token:          token,
				expectedStatus: http.StatusNoContent,
			},
		}

		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/categories/%s", tc.id), nil)
				req.Header.Set("Authorization", "Bearer "+tc.token)
				w := httptest.NewRecorder()

				router.ServeHTTP(w, req)

				if tc.expectedStatus == http.StatusNotFound {
					assert.Contains(t, []int{http.StatusNotFound, http.StatusForbidden}, w.Code)
				} else {
					assert.Equal(t, tc.expectedStatus, w.Code)
				}

				if tc.expectedStatus == http.StatusNoContent {
					var count int
					err := res.DB.QueryRow("SELECT COUNT(*) FROM categories WHERE id = $1", tc.id).Scan(&count)
					assert.NoError(t, err)
					assert.Equal(t, 0, count)
				} else if tc.id == categoryID {
					var count int
					err := res.DB.QueryRow("SELECT COUNT(*) FROM categories WHERE id = $1", tc.id).Scan(&count)
					assert.NoError(t, err)
					assert.Equal(t, 1, count, "Category should still exist after failed delete")
				}
			})
		}
	})
}
