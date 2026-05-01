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
	"github.com/stretchr/testify/assert"

	"github.com/felipe1496/open-wallet/infra"
	"github.com/felipe1496/open-wallet/internal/factory"
	"github.com/felipe1496/open-wallet/internal/resources/recurrences/handlers"
	"github.com/felipe1496/open-wallet/internal/routes"
	"github.com/felipe1496/open-wallet/internal/utils"
)

func setupRecurrenceTestServer(pg *sql.DB, db *sql.DB, cfg *infra.Config) (*gin.Engine, *factory.Factory) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	f := factory.NewFactory(pg, cfg)

	routes.SetupRecurrencesRoutes(router, f, cfg)

	return router, f
}

func TestE2eRecurrences(t *testing.T) {
	res := SetupTestResources(t)
	defer func() { _ = res.PostgresContainer.Terminate(context.Background()) }()
	defer func() { _ = res.DB.Close() }()

	cfg := &infra.Config{
		JWTSecret: "test-secret",
	}
	cfg.RateLimits.XS = func() (int, int) { return 1000, 60000 }
	cfg.RateLimits.SM = func() (int, int) { return 1000, 60000 }
	cfg.RateLimits.MD = func() (int, int) { return 1000, 60000 }

	AssertTableIsEmpty(t, res.DB, "users")
	AssertTableIsEmpty(t, res.DB, "recurrences")

	testUser, token := SetupTestUser(t, res.DB, cfg)
	router, _ := setupRecurrenceTestServer(res.DB, res.DB, cfg)

	t.Run("Authentication Enforcement", func(t *testing.T) {
		endpoints := []struct {
			method string
			url    string
		}{
			{http.MethodPost, "/api/v1/recurrences"},
			{http.MethodGet, "/api/v1/recurrences"},
			{http.MethodPost, "/api/v1/recurrences/202604"},
			{http.MethodPatch, "/api/v1/recurrences/some-id"},
			{http.MethodDelete, "/api/v1/recurrences/some-id"},
		}

		for _, e := range endpoints {
			AssertUnauthorized(t, router, e.method, e.url, nil)
		}
	})

	t.Run("POST /recurrences", func(t *testing.T) {
		endPeriod := "202612"

		type testCase struct {
			name            string
			payload         handlers.CreateRecurrenceRequest
			expectedStatus  int
			validateDB      bool
			expectEndPeriod *string
		}

		cases := []testCase{
			{
				name: "should create recurrence successfully (only start_period)",
				payload: handlers.CreateRecurrenceRequest{
					Name:        "Internet Bill",
					Amount:      -100.0,
					DayOfMonth:  15,
					StartPeriod: "202604",
				},
				expectedStatus:  http.StatusCreated,
				validateDB:      true,
				expectEndPeriod: nil,
			},
			{
				name: "should create recurrence successfully (start_period and end_period)",
				payload: handlers.CreateRecurrenceRequest{
					Name:        "Car Lease",
					Amount:      -300.0,
					DayOfMonth:  5,
					StartPeriod: "202601",
					EndPeriod:   &endPeriod,
				},
				expectedStatus:  http.StatusCreated,
				validateDB:      true,
				expectEndPeriod: &endPeriod,
			},
			{
				name: "should fail with invalid payload",
				payload: handlers.CreateRecurrenceRequest{
					Name: "", // invalid
				},
				expectedStatus: http.StatusBadRequest,
				validateDB:     false,
			},
		}

		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				body, _ := json.Marshal(tc.payload)
				req := httptest.NewRequest(http.MethodPost, "/api/v1/recurrences", bytes.NewBuffer(body))
				req.Header.Set("Authorization", "Bearer "+token)
				req.Header.Set("Content-Type", "application/json")
				w := httptest.NewRecorder()

				router.ServeHTTP(w, req)

				assert.Equal(t, tc.expectedStatus, w.Code)

				if tc.validateDB {
					var count int
					var storedEndPeriod *string
					err := res.DB.QueryRow("SELECT COUNT(*), MAX(end_period) FROM recurrences WHERE name = $1 AND user_id = $2", tc.payload.Name, testUser.ID).Scan(&count, &storedEndPeriod)
					assert.NoError(t, err)
					assert.Equal(t, 1, count)
					if tc.expectEndPeriod == nil {
						assert.Nil(t, storedEndPeriod)
					} else {
						assert.NotNil(t, storedEndPeriod)
						assert.Equal(t, *tc.expectEndPeriod, *storedEndPeriod)
					}
				}
			})
		}
	})

	t.Run("GET /recurrences", func(t *testing.T) {
		// Clear and Seed
		_, _ = res.DB.Exec("DELETE FROM recurrences")
		_, err := res.DB.Exec("INSERT INTO recurrences (id, user_id, name, amount, day_of_month, start_period) VALUES ($1, $2, $3, $4, $5, $6)",
			ulid.Make().String(), testUser.ID, "Rent", -1000.0, 1, "202601")
		assert.NoError(t, err)

		req := httptest.NewRequest(http.MethodGet, "/api/v1/recurrences", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response utils.PaginatedResponse[handlers.ListRecurrencesResponseData]
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Len(t, response.Data.Recurrences, 1)
		assert.Equal(t, "Rent", response.Data.Recurrences[0].Name)
	})

	t.Run("GET /recurrences - Filtering", func(t *testing.T) {
		type testCase struct {
			name          string
			queryString   string
			expectedCount int
		}

		cases := []testCase{
			{
				name:          "filter by name exact",
				queryString:   "filter=name eq 'Rent'",
				expectedCount: 1,
			},
			{
				name:          "filter by amount",
				queryString:   "filter=amount eq -1000",
				expectedCount: 1,
			},
			{
				name:          "no match",
				queryString:   "filter=name eq 'Gym'",
				expectedCount: 0,
			},
		}

		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				params := url.Values{}
				params.Set("filter", tc.queryString[7:]) // Strip "filter=" prefix

				req := httptest.NewRequest(http.MethodGet, "/api/v1/recurrences?"+params.Encode(), nil)
				req.Header.Set("Authorization", "Bearer "+token)
				w := httptest.NewRecorder()

				router.ServeHTTP(w, req)

				assert.Equal(t, http.StatusOK, w.Code)
				var response utils.PaginatedResponse[handlers.ListRecurrencesResponseData]
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Len(t, response.Data.Recurrences, tc.expectedCount)
			})
		}
	})

	t.Run("POST /recurrences/:period (Prepare)", func(t *testing.T) {
		// Setup Category for the recurrence
		categoryID := ulid.Make().String()
		_, _ = res.DB.Exec("INSERT INTO categories (id, user_id, name, color) VALUES ($1, $2, $3, $4)",
			categoryID, testUser.ID, "Fixed", "#000000")

		// Create a recurrence template with start_period 202604 but no end_period
		recNoEnd := ulid.Make().String()
		_, _ = res.DB.Exec("INSERT INTO recurrences (id, user_id, name, amount, day_of_month, start_period, category_id) VALUES ($1, $2, $3, $4, $5, $6, $7)",
			recNoEnd, testUser.ID, "No End Rec", -50.0, 10, "202604", categoryID)

		// Create a recurrence template with start_period 202604 and end_period 202606
		recWithEnd := ulid.Make().String()
		_, _ = res.DB.Exec("INSERT INTO recurrences (id, user_id, name, amount, day_of_month, start_period, end_period, category_id) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)",
			recWithEnd, testUser.ID, "With End Rec", -100.0, 15, "202604", "202606", categoryID)

		tr, f := setupRecurrenceTestServer(res.DB, res.DB, cfg)
		routes.SetupTransactionsRoutes(tr, f, cfg)

		// Helper to query transaction count for a recurrence in a period
		getTxCount := func(recID, period string) int {
			var count int
			err := res.DB.QueryRow(`
				SELECT COUNT(t.id) 
				FROM transactions t 
				JOIN entries e ON e.transaction_id = t.id 
				WHERE t.recurrence_id = $1 AND TO_CHAR(e.reference_date, 'YYYYMM') = $2`, recID, period).Scan(&count)
			assert.NoError(t, err)
			return count
		}

		type prepareCase struct {
			period          string
			expectedNoEnd   int
			expectedWithEnd int
		}

		// Test multiple periods
		cases := []prepareCase{
			{"202603", 0, 0}, // Before start_period
			{"202604", 1, 1}, // At start_period
			{"202605", 1, 1}, // Inside period
			{"202606", 1, 1}, // At end_period
			{"202607", 1, 0}, // After end_period for recWithEnd
		}

		for _, tc := range cases {
			t.Run("Prepare for period "+tc.period, func(t *testing.T) {
				req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/v1/recurrences/%s", tc.period), nil)
				req.Header.Set("Authorization", "Bearer "+token)
				w := httptest.NewRecorder()

				tr.ServeHTTP(w, req)

				assert.Equal(t, http.StatusNoContent, w.Code)

				// Verify transactions were generated correctly
				assert.Equal(t, tc.expectedNoEnd, getTxCount(recNoEnd, tc.period), "Mismatch for recNoEnd in period "+tc.period)
				assert.Equal(t, tc.expectedWithEnd, getTxCount(recWithEnd, tc.period), "Mismatch for recWithEnd in period "+tc.period)
			})
		}
	})

	t.Run("DELETE /recurrences/:id", func(t *testing.T) {
		recurrenceID := ulid.Make().String()
		_, _ = res.DB.Exec("INSERT INTO recurrences (id, user_id, name, amount, day_of_month, start_period) VALUES ($1, $2, $3, $4, $5, $6)",
			recurrenceID, testUser.ID, "To Delete", -10.0, 1, "202604")

		_, tokenB := SetupTestUser(t, res.DB, cfg)

		type testCase struct {
			name           string
			id             string
			token          string
			expectedStatus int
		}

		cases := []testCase{
			{
				name:           "should fail when deleting non-existent recurrence",
				id:             ulid.Make().String(),
				token:          token,
				expectedStatus: http.StatusNotFound,
			},
			{
				name:           "should fail when deleting recurrence of another user",
				id:             recurrenceID,
				token:          tokenB,
				expectedStatus: http.StatusNotFound,
			},
			{
				name:           "should delete successfully",
				id:             recurrenceID,
				token:          token,
				expectedStatus: http.StatusNoContent,
			},
		}

		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/recurrences/%s?scope=all", tc.id), nil)
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
					_ = res.DB.QueryRow("SELECT COUNT(*) FROM recurrences WHERE id = $1", tc.id).Scan(&count)
					assert.Equal(t, 0, count)
				} else if tc.id == recurrenceID {
					var count int
					_ = res.DB.QueryRow("SELECT COUNT(*) FROM recurrences WHERE id = $1", tc.id).Scan(&count)
					assert.Equal(t, 1, count, "Recurrence should still exist after failed delete")
				}
			})
		}
	})
}
