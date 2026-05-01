package e2e

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/felipe1496/open-wallet/internal/resources/auth/handlers"
	authUseCases "github.com/felipe1496/open-wallet/internal/resources/auth/usecases"
	"github.com/felipe1496/open-wallet/internal/resources/users/repository"
	usersUseCases "github.com/felipe1496/open-wallet/internal/resources/users/usecases"
	"github.com/felipe1496/open-wallet/internal/services"
	"github.com/felipe1496/open-wallet/internal/services/mocks"
)

func setupTestServer(db *sql.DB, googleService services.GoogleService, jwtService services.JWTService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	usersRepo := repository.NewUsersRepo()
	usersUseCase := usersUseCases.NewUsersUseCases(usersRepo, db)
	authUseCase := authUseCases.NewAuthUseCases(googleService, usersUseCase)
	handler := handlers.NewHandler(authUseCase, jwtService)
	authGroup := router.Group("/api/v1/auth")
	{
		authGroup.POST("/login/google", handler.CreateLoginWithGoogle)
	}

	return router
}

func TestE2eAuth(t *testing.T) {
	res := SetupTestResources(t)
	defer func() { _ = res.PostgresContainer.Terminate(context.Background()) }()
	defer func() { _ = res.DB.Close() }()

	t.Log("postgres started with:", res.PostgresConnStr)

	AssertTableIsEmpty(t, res.DB, "users")

	t.Run("should create user when login with Google for the first time", func(t *testing.T) {
		mockGoogleService := new(mocks.MockGoogleService)
		mockJWTService := new(mocks.MockJWTService)

		accessToken := "mock-access-token"
		mockGoogleService.
			On("GetUserAccessToken", "valid-code").
			Return(&accessToken, nil)

		email := "test@example.com"
		emailVerified := true
		picture := "https://example.com/avatar.jpg"
		mockGoogleService.On("GetUserAccessToken", "valid-code").Return(&accessToken, nil)
		mockGoogleService.
			On("GetUserInfo", accessToken).
			Return(&services.GoogleUserInfo{
				Sub:           "google-sub-123",
				Name:          "Test User",
				Email:         &email,
				EmailVerified: &emailVerified,
				Picture:       &picture,
			}, nil)
		mockJWTService.On("GenerateToken", mock.Anything).Return("mock-token", nil)

		router := setupTestServer(res.DB, mockGoogleService, mockJWTService)

		body := handlers.LoginGoogleRequest{Code: "valid-code"}
		bodyJSON, _ := json.Marshal(body)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login/google", bytes.NewBuffer(bodyJSON))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]map[string]any
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response["data"])
		assert.NotNil(t, response["data"]["user"])
		assert.NotEmpty(t, response["data"]["access_token"])

		var count int
		err = res.DB.QueryRow("SELECT COUNT(*) FROM users WHERE email = $1", email).Scan(&count)
		assert.NoError(t, err)
		assert.Equal(t, 1, count, "should have exactly one user with this email")

		var dbUser struct {
			ID        string
			Name      string
			Email     string
			AvatarURL string
			Username  string
		}
		err = res.DB.QueryRow(
			"SELECT id, name, email, avatar_url, username FROM users WHERE email = $1",
			email,
		).Scan(&dbUser.ID, &dbUser.Name, &dbUser.Email, &dbUser.AvatarURL, &dbUser.Username)

		assert.NoError(t, err)
		assert.NotEmpty(t, dbUser.ID)
		assert.Equal(t, "Test User", dbUser.Name)
		assert.Equal(t, email, dbUser.Email)
		assert.Equal(t, picture, dbUser.AvatarURL)
		assert.NotEmpty(t, dbUser.Username)

		mockGoogleService.AssertExpectations(t)
	})

	t.Run("should return existing user when login with Google for second time", func(t *testing.T) {
		email := "test@example.com"
		var count int
		err := res.DB.QueryRow("SELECT COUNT(*) FROM users WHERE email = $1", email).Scan(&count)
		assert.NoError(t, err)
		assert.Equal(t, 1, count, "user should already exist before testing second login")

		mockGoogleService := new(mocks.MockGoogleService)
		mockJWTService := new(mocks.MockJWTService)

		accessToken := "mock-access-token"
		mockGoogleService.
			On("GetUserAccessToken", "valid-code").
			Return(&accessToken, nil)

		emailVerified := true
		picture := "https://example.com/avatar.jpg"
		mockGoogleService.
			On("GetUserInfo", accessToken).
			Return(&services.GoogleUserInfo{
				Sub:           "google-sub-123",
				Name:          "Test User",
				Email:         &email,
				EmailVerified: &emailVerified,
				Picture:       &picture,
			}, nil)

		mockJWTService.On("GenerateToken", mock.Anything).Return("mock-token", nil)

		router := setupTestServer(res.DB, mockGoogleService, mockJWTService)

		body := handlers.LoginGoogleRequest{Code: "valid-code"}
		bodyJSON, _ := json.Marshal(body)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login/google", bytes.NewBuffer(bodyJSON))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response map[string]map[string]any
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		assert.NotNil(t, response["data"])
		assert.NotNil(t, response["data"]["user"])
		assert.NotEmpty(t, response["data"]["access_token"])

		err = res.DB.QueryRow("SELECT COUNT(*) FROM users WHERE email = $1", email).Scan(&count)
		assert.NoError(t, err)
		assert.Equal(t, 1, count, "should have exactly one user with this email")

		var dbUser struct {
			ID        string
			Name      string
			Email     string
			AvatarURL string
			Username  string
		}
		err = res.DB.QueryRow(
			"SELECT id, name, email, avatar_url, username FROM users WHERE email = $1",
			email,
		).Scan(&dbUser.ID, &dbUser.Name, &dbUser.Email, &dbUser.AvatarURL, &dbUser.Username)

		assert.NoError(t, err)
		assert.NotEmpty(t, dbUser.ID)
		assert.Equal(t, "Test User", dbUser.Name)
		assert.Equal(t, email, dbUser.Email)
		assert.Equal(t, picture, dbUser.AvatarURL)
		assert.NotEmpty(t, dbUser.Username)

		mockGoogleService.AssertExpectations(t)
	})

	// should reject invalid code and not create a new user
	t.Run("should reject invalid code and not create user", func(t *testing.T) {
		mockGoogleService := new(mocks.MockGoogleService)
		mockJWTService := new(mocks.MockJWTService)

		emptyAccessToken := ""
		mockGoogleService.
			On("GetUserAccessToken", "invalid-code").
			Return(&emptyAccessToken, services.FailedGoogleAuthenticationErr)

		router := setupTestServer(res.DB, mockGoogleService, mockJWTService)

		body := handlers.LoginGoogleRequest{Code: "invalid-code"}
		bodyJSON, _ := json.Marshal(body)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login/google", bytes.NewBuffer(bodyJSON))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)

		var count int
		err := res.DB.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
		assert.NoError(t, err)
		assert.Equal(t, 1, count, "no new user should be created on invalid code")

		mockGoogleService.AssertExpectations(t)
	})
}
