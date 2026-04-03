package auth

import "github.com/felipe1496/open-wallet/internal/resources/users"

// ==============================================================================
// 1. HTTP MODELS
//    Models that represents request or response objects
// ==============================================================================

type LoginGoogleRequest struct {
	Code string `json:"code"`
}

type LoginGoogleResponse struct {
	Data LoginGoogleResponseData `json:"data"`
}

type LoginGoogleResponseData struct {
	AccessToken string     `json:"access_token"`
	User        users.User `json:"user"`
}
