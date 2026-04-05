package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/felipe1496/open-wallet/infra"
)

type GoogleService interface {
	GetUserInfo(accessToken string) (*GoogleUserInfo, error)
	GetUserAccessToken(code string) (*string, error)
}

func NewGoogleService(cfg *infra.Config) GoogleService {
	return &googleServiceImpl{
		cfg: cfg,
	}
}

type googleServiceImpl struct {
	cfg *infra.Config
}

type GoogleUserInfo struct {
	Sub           string  `json:"sub"`
	Name          string  `json:"name"`
	GivenName     *string `json:"given_name"`
	FamilyName    *string `json:"family_name"`
	Picture       *string `json:"picture"`
	Email         *string `json:"email"`
	EmailVerified *bool   `json:"email_verified"`
}

func (s *googleServiceImpl) GetUserInfo(accessToken string) (*GoogleUserInfo, error) {
	req, _ := http.NewRequest("GET", "https://www.googleapis.com/oauth2/v3/userinfo", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return nil, FailedGoogleAuthenticationErr
	}
	defer func() { _ = response.Body.Close() }()

	if response.StatusCode != http.StatusOK {
		return nil, FailedGoogleAuthenticationErr
	}

	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, FailedGoogleAuthenticationErr
	}

	var user GoogleUserInfo
	if err := json.Unmarshal(bodyBytes, &user); err != nil {
		return nil, FailedGoogleAuthenticationErr
	}

	return &user, nil
}

func (s *googleServiceImpl) GetUserAccessToken(code string) (*string, error) {
	data := url.Values{}
	data.Set("code", code)
	data.Set("client_id", s.cfg.GoogleClientID)
	data.Set("client_secret", s.cfg.GoogleSecret)
	data.Set("redirect_uri", s.cfg.LoginRedirectURI)
	data.Set("grant_type", "authorization_code")

	req, err := http.NewRequest(
		"POST",
		"https://oauth2.googleapis.com/token",
		strings.NewReader(data.Encode()),
	)
	if err != nil {
		return nil, FailedGoogleAuthenticationErr
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, FailedGoogleAuthenticationErr
	}

	if res.StatusCode != http.StatusOK {
		return nil, FailedGoogleAuthenticationErr
	}

	defer func() { _ = res.Body.Close() }()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, FailedGoogleAuthenticationErr
	}

	var response struct {
		AccessToken string `json:"access_token"`
	}

	if err := json.Unmarshal(body, &response); err != nil {
		return nil, FailedGoogleAuthenticationErr
	}

	return &response.AccessToken, nil
}
