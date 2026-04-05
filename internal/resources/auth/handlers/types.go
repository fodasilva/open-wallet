package handlers

import (
	"github.com/felipe1496/open-wallet/internal/resources/users/repository"
)

type LoginGoogleRequest struct {
	Code string `json:"code"`
}

type LoginGoogleResponse struct {
	Data LoginGoogleResponseData `json:"data"`
}

type LoginGoogleResponseData struct {
	AccessToken string          `json:"access_token"`
	User        repository.User `json:"user"`
}
