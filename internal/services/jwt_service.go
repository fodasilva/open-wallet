package services

import (
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/felipe1496/open-wallet/infra"
	"github.com/felipe1496/open-wallet/internal/utils"
)

type JWTService interface {
	GenerateToken(userId string) (string, error)
	ValidateToken(tokenString string) (string, error)
}

type JWTServiceImpl struct {
	cfg *infra.Config
}

func NewJWTService(cfg *infra.Config) JWTService {
	return &JWTServiceImpl{
		cfg: cfg,
	}
}

func (s *JWTServiceImpl) GenerateToken(userID string) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
		"iat": time.Now().Unix(),
		"iss": "money-api",
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(s.cfg.JWTSecret))

	if err != nil {
		return "", utils.NewHTTPError(http.StatusInternalServerError, "failed to generate JWT token")
	}

	return signedToken, nil
}

func (s *JWTServiceImpl) ValidateToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("método de assinatura inválido: %v", token.Header["alg"])
		}
		return []byte(s.cfg.JWTSecret), nil
	})

	if err != nil {
		return "", utils.NewHTTPError(http.StatusInternalServerError, "failed to parse JWT token")
	}

	if !token.Valid {
		return "", utils.NewHTTPError(http.StatusUnauthorized, "invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", utils.NewHTTPError(http.StatusInternalServerError, "failed to extract claims from token")
	}

	if iss, ok := claims["iss"].(string); !ok || iss != "money-api" {
		return "", utils.NewHTTPError(http.StatusUnauthorized, "invalid issuer")
	}

	userID, ok := claims["sub"].(string)
	if !ok {
		return "", utils.NewHTTPError(http.StatusInternalServerError, "failed to extract sub claim from token")
	}

	return userID, nil
}
