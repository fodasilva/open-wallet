package e2e

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/oklog/ulid/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"

	"github.com/felipe1496/open-wallet/infra"
	"github.com/felipe1496/open-wallet/internal/factory"
	"github.com/felipe1496/open-wallet/internal/resources/transactions/handlers"
	"github.com/felipe1496/open-wallet/internal/routes"
	"github.com/felipe1496/open-wallet/internal/utils"
)

func setupTransactionTestServer(pg *sql.DB, redisClient *redis.Client, cfg *infra.Config) (*gin.Engine, *factory.Factory) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	f := factory.NewFactory(pg, cfg)

	routes.SetupTransactionsRoutes(router, f, redisClient, cfg)

	return router, f
}

func TestE2eTransactions(t *testing.T) {
	res := SetupTestResources(t)
	defer func() { _ = res.PostgresContainer.Terminate(context.Background()) }()
	defer func() { _ = res.RedisContainer.Terminate(context.Background()) }()
	defer func() { _ = res.DB.Close() }()

	cfg := &infra.Config{
		JWTSecret: "test-secret",
	}
	cfg.RateLimits.XS = func() (int, int) { return 1000, 60000 }

	AssertTableIsEmpty(t, res.DB, "users")
	AssertTableIsEmpty(t, res.DB, "transactions")
	AssertTableIsEmpty(t, res.DB, "entries")

	testUser, token := SetupTestUser(t, res.DB, cfg)
	router, _ := setupTransactionTestServer(res.DB, res.RedisClient, cfg)

	t.Run("Authentication Enforcement", func(t *testing.T) {
		endpoints := []struct {
			method string
			url    string
		}{
			{http.MethodPost, "/api/v1/transactions"},
			{http.MethodGet, "/api/v1/transactions/entries"},
			{http.MethodPatch, "/api/v1/transactions/some-id"},
			{http.MethodDelete, "/api/v1/transactions/some-id"},
		}

		for _, e := range endpoints {
			AssertUnauthorized(t, router, e.method, e.url, nil)
		}
	})

	t.Run("POST /transactions", func(t *testing.T) {
		// Seed a category
		categoryID := ulid.Make().String()
		_, err := res.DB.Exec("INSERT INTO categories (id, user_id, name, color) VALUES ($1, $2, $3, $4)",
			categoryID, testUser.ID, "Food", "#FF0000")
		assert.NoError(t, err)

		type testCase struct {
			name           string
			payload        handlers.CreateTransactionRequest
			expectedStatus int
			validateDB     bool
		}

		cases := []testCase{
			{
				name: "should create transaction successfully",
				payload: handlers.CreateTransactionRequest{
					Name:       "Lunch",
					CategoryID: &categoryID,
					Type:       "simple_expense",
					Entries: []handlers.CreateEntryRequest{
						{Amount: -25.5, ReferenceDate: "2026-04-10"},
					},
				},
				expectedStatus: http.StatusCreated,
				validateDB:     true,
			},
			{
				name: "should fail with invalid payload (missing entries)",
				payload: handlers.CreateTransactionRequest{
					Name: "Invalid",
					Type: "simple_expense",
				},
				expectedStatus: http.StatusBadRequest,
				validateDB:     false,
			},
		}

		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				body, _ := json.Marshal(tc.payload)
				req := httptest.NewRequest(http.MethodPost, "/api/v1/transactions", bytes.NewBuffer(body))
				req.Header.Set("Authorization", "Bearer "+token)
				req.Header.Set("Content-Type", "application/json")
				w := httptest.NewRecorder()

				router.ServeHTTP(w, req)

				assert.Equal(t, tc.expectedStatus, w.Code)

				if tc.validateDB {
					var count int
					err := res.DB.QueryRow("SELECT COUNT(*) FROM transactions WHERE name = $1 AND user_id = $2", tc.payload.Name, testUser.ID).Scan(&count)
					assert.NoError(t, err)
					assert.Equal(t, 1, count)

					var entryCount int
					err = res.DB.QueryRow("SELECT COUNT(*) FROM entries e JOIN transactions t ON t.id = e.transaction_id WHERE t.name = $1", tc.payload.Name).Scan(&entryCount)
					assert.NoError(t, err)
					assert.Equal(t, len(tc.payload.Entries), entryCount)
				}
			})
		}
	})

	t.Run("GET /transactions/entries", func(t *testing.T) {
		type testCase struct {
			name           string
			seedFunc       func()
			expectedCount  int
			expectedStatus int
		}

		cases := []testCase{
			{
				name: "should list entries successfully",
				seedFunc: func() {
					_, _ = res.DB.Exec("DELETE FROM entries")
					_, _ = res.DB.Exec("DELETE FROM transactions")
					txID := ulid.Make().String()
					_, _ = res.DB.Exec("INSERT INTO transactions (id, user_id, name, category) VALUES ($1, $2, $3, $4)",
						txID, testUser.ID, "Salary", "income")
					_, _ = res.DB.Exec("INSERT INTO entries (id, transaction_id, amount, reference_date) VALUES ($1, $2, $3, $4)",
						ulid.Make().String(), txID, 5000.0, "2026-04-01")
				},
				expectedCount:  1,
				expectedStatus: http.StatusOK,
			},
		}

		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				tc.seedFunc()

				req := httptest.NewRequest(http.MethodGet, "/api/v1/transactions/entries", nil)
				req.Header.Set("Authorization", "Bearer "+token)
				w := httptest.NewRecorder()

				router.ServeHTTP(w, req)

				assert.Equal(t, tc.expectedStatus, w.Code)
				if tc.expectedStatus == http.StatusOK {
					var response utils.PaginatedResponse[handlers.ListEntriesResponseData]
					err := json.Unmarshal(w.Body.Bytes(), &response)
					assert.NoError(t, err)
					assert.Len(t, response.Data.Entries, tc.expectedCount)
				}
			})
		}
	})

	t.Run("GET /transactions/entries - Filtering", func(t *testing.T) {
		type testCase struct {
			name          string
			queryString   string
			expectedCount int
		}

		cases := []testCase{
			{
				name:          "filter by amount exact",
				queryString:   "filter=amount eq 5000",
				expectedCount: 1,
			},
			{
				name:          "filter by amount range",
				queryString:   "filter=amount gt 1000",
				expectedCount: 1,
			},
			{
				name:          "filter by reference date",
				queryString:   "filter=reference_date eq '2026-04-01'",
				expectedCount: 1,
			},
			{
				name:          "no match",
				queryString:   "filter=amount lt 0",
				expectedCount: 0,
			},
		}

		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				params := url.Values{}
				params.Set("filter", tc.queryString[7:]) // Strip "filter=" prefix

				req := httptest.NewRequest(http.MethodGet, "/api/v1/transactions/entries?"+params.Encode(), nil)
				req.Header.Set("Authorization", "Bearer "+token)
				w := httptest.NewRecorder()

				router.ServeHTTP(w, req)

				assert.Equal(t, http.StatusOK, w.Code)
				var response utils.PaginatedResponse[handlers.ListEntriesResponseData]
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Len(t, response.Data.Entries, tc.expectedCount)
			})
		}
	})

	t.Run("PATCH /transactions/:id", func(t *testing.T) {
		_, _ = res.DB.Exec("DELETE FROM entries")
		_, _ = res.DB.Exec("DELETE FROM transactions")

		txID := ulid.Make().String()
		_, _ = res.DB.Exec("INSERT INTO transactions (id, user_id, name, category) VALUES ($1, $2, $3, $4)",
			txID, testUser.ID, "Old Name", "simple_expense")
		_, _ = res.DB.Exec("INSERT INTO entries (id, transaction_id, amount, reference_date) VALUES ($1, $2, $3, $4)",
			ulid.Make().String(), txID, -50.0, "2026-04-10")

		type testCase struct {
			name           string
			id             string
			payload        handlers.UpdateTransactionRequest
			expectedStatus int
			expectedName   string
		}

		newName := "New Name"
		cases := []testCase{
			{
				name: "should update name successfully",
				id:   txID,
				payload: handlers.UpdateTransactionRequest{
					Name: &newName,
				},
				expectedStatus: http.StatusOK,
				expectedName:   newName,
			},
			{
				name:           "should fail with 404 for non-existent transaction",
				id:             ulid.Make().String(),
				payload:        handlers.UpdateTransactionRequest{Name: &newName},
				expectedStatus: http.StatusNotFound,
			},
		}

		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				body, _ := json.Marshal(tc.payload)
				req := httptest.NewRequest(http.MethodPatch, fmt.Sprintf("/api/v1/transactions/%s", tc.id), bytes.NewBuffer(body))
				req.Header.Set("Authorization", "Bearer "+token)
				req.Header.Set("Content-Type", "application/json")
				w := httptest.NewRecorder()

				router.ServeHTTP(w, req)

				assert.Equal(t, tc.expectedStatus, w.Code)

				if tc.expectedStatus == http.StatusOK {
					var storedName string
					err := res.DB.QueryRow("SELECT name FROM transactions WHERE id = $1", tc.id).Scan(&storedName)
					assert.NoError(t, err)
					assert.Equal(t, tc.expectedName, storedName)
				}
			})
		}
	})

	t.Run("DELETE /transactions/:id", func(t *testing.T) {
		type testCase struct {
			name           string
			seedFunc       func() string
			expectedStatus int
		}

		cases := []testCase{
			{
				name: "should delete successfully",
				seedFunc: func() string {
					_, _ = res.DB.Exec("DELETE FROM entries")
					_, _ = res.DB.Exec("DELETE FROM transactions")
					txID := ulid.Make().String()
					_, _ = res.DB.Exec("INSERT INTO transactions (id, user_id, name, category) VALUES ($1, $2, $3, $4)",
						txID, testUser.ID, "To Delete", "simple_expense")
					_, _ = res.DB.Exec("INSERT INTO entries (id, transaction_id, amount, reference_date) VALUES ($1, $2, $3, $4)",
						ulid.Make().String(), txID, -10.0, "2026-04-10")
					return txID
				},
				expectedStatus: http.StatusNoContent,
			},
			{
				name: "should return 404 when transaction is not found",
				seedFunc: func() string {
					return ulid.Make().String()
				},
				expectedStatus: http.StatusNotFound,
			},
		}

		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				id := tc.seedFunc()

				req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/transactions/%s", id), nil)
				req.Header.Set("Authorization", "Bearer "+token)
				w := httptest.NewRecorder()

				router.ServeHTTP(w, req)

				assert.Equal(t, tc.expectedStatus, w.Code)

				if tc.expectedStatus == http.StatusNoContent {
					var count int
					_ = res.DB.QueryRow("SELECT COUNT(*) FROM transactions WHERE id = $1", id).Scan(&count)
					assert.Equal(t, 0, count)
				}
			})
		}
	})

	t.Run("GET /summary", func(t *testing.T) {
		// Seed transactions for multiple periods
		// Period: 202601
		_, _ = res.DB.Exec("INSERT INTO transactions (id, user_id, name, category) VALUES ($1, $2, $3, $4)",
			ulid.Make().String(), testUser.ID, "Jan Expense", "simple_expense")
		_, _ = res.DB.Exec("INSERT INTO entries (id, transaction_id, amount, reference_date) VALUES ($1, (SELECT id FROM transactions WHERE name = 'Jan Expense'), $2, $3)",
			ulid.Make().String(), 100.0, "2026-01-15")

		_, _ = res.DB.Exec("INSERT INTO transactions (id, user_id, name, category) VALUES ($1, $2, $3, $4)",
			ulid.Make().String(), testUser.ID, "Jan Income", "income")
		_, _ = res.DB.Exec("INSERT INTO entries (id, transaction_id, amount, reference_date) VALUES ($1, (SELECT id FROM transactions WHERE name = 'Jan Income'), $2, $3)",
			ulid.Make().String(), 500.0, "2026-01-10")

		// Period: 202602
		_, _ = res.DB.Exec("INSERT INTO transactions (id, user_id, name, category) VALUES ($1, $2, $3, $4)",
			ulid.Make().String(), testUser.ID, "Feb Expense", "simple_expense")
		_, _ = res.DB.Exec("INSERT INTO entries (id, transaction_id, amount, reference_date) VALUES ($1, (SELECT id FROM transactions WHERE name = 'Feb Expense'), $2, $3)",
			ulid.Make().String(), 250.0, "2026-02-20")

		t.Run("should fail when missing period filters", func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/api/v1/transactions/summary", nil)
			req.Header.Set("Authorization", "Bearer "+token)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusBadRequest, w.Code)
		})

		t.Run("should fail when period format is invalid", func(t *testing.T) {
			cases := []struct {
				name   string
				filter string
			}{
				{"too short", "period gte '20261' and period lte '202602'"},
				{"too long", "period gte '2026011' and period lte '202602'"},
				{"invalid month", "period gte '202613' and period lte '202602'"},
				{"non-digits", "period gte '2026AB' and period lte '202602'"},
				{"gte > lte", "period gte '202602' and period lte '202601'"},
				{"range > 12 months", "period gte '202601' and period lte '202701'"},
			}

			for _, tc := range cases {
				t.Run(tc.name, func(t *testing.T) {
					params := url.Values{}
					params.Set("filter", tc.filter)

					req := httptest.NewRequest(http.MethodGet, "/api/v1/transactions/summary?"+params.Encode(), nil)
					req.Header.Set("Authorization", "Bearer "+token)
					w := httptest.NewRecorder()

					router.ServeHTTP(w, req)

					assert.Equal(t, http.StatusBadRequest, w.Code)
				})
			}
		})

		t.Run("should return summary for valid period range", func(t *testing.T) {
			params := url.Values{}
			params.Set("filter", "period gte '202601' and period lte '202602'")

			req := httptest.NewRequest(http.MethodGet, "/api/v1/transactions/summary?"+params.Encode(), nil)
			req.Header.Set("Authorization", "Bearer "+token)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)

			var response utils.ResponseData[handlers.SummaryResponseData]
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)

			assert.Len(t, response.Data.Summary, 2)

			// Check 202601
			foundJan := false
			for _, s := range response.Data.Summary {
				if s.Period == "202601" {
					assert.Equal(t, 500.0, s.Income)
					assert.Equal(t, 100.0, s.Expense)
					assert.Equal(t, 400.0, s.Balance)
					foundJan = true
				}
			}
			assert.True(t, foundJan)
		})

		t.Run("should respect sorting", func(t *testing.T) {
			params := url.Values{}
			params.Set("filter", "period gte '202601' and period lte '202602'")
			params.Set("order_by", "period:desc")

			req := httptest.NewRequest(http.MethodGet, "/api/v1/transactions/summary?"+params.Encode(), nil)
			req.Header.Set("Authorization", "Bearer "+token)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)

			var response utils.ResponseData[handlers.SummaryResponseData]
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)

			assert.Len(t, response.Data.Summary, 2)
			assert.Equal(t, "202602", response.Data.Summary[0].Period)
			assert.Equal(t, "202601", response.Data.Summary[1].Period)
		})
		t.Run("should skip periods with no data", func(t *testing.T) {
			// Seed data for 202604
			_, _ = res.DB.Exec("INSERT INTO transactions (id, user_id, name, category) VALUES ($1, $2, $3, $4)",
				ulid.Make().String(), testUser.ID, "April Expense", "simple_expense")
			_, _ = res.DB.Exec("INSERT INTO entries (id, transaction_id, amount, reference_date) VALUES ($1, (SELECT id FROM transactions WHERE name = 'April Expense'), $2, $3)",
				ulid.Make().String(), 50.0, "2026-04-05")

			// Request range 202602 to 202604. 202603 is empty.
			params := url.Values{}
			params.Set("filter", "period gte '202602' and period lte '202604'")
			params.Set("order_by", "period:asc")

			req := httptest.NewRequest(http.MethodGet, "/api/v1/transactions/summary?"+params.Encode(), nil)
			req.Header.Set("Authorization", "Bearer "+token)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)

			var response utils.ResponseData[handlers.SummaryResponseData]
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)

			// Should only have 202602 and 202604. 202603 should NOT be there.
			assert.Len(t, response.Data.Summary, 2)
			assert.Equal(t, "202602", response.Data.Summary[0].Period)
			assert.Equal(t, "202604", response.Data.Summary[1].Period)
		})
	})
}
