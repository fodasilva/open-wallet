package handlers

import (
	"time"
)

type LoginGoogleRequest struct {
	Code string `json:"code"`
}

type UserResource struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	AvatarURL string    `json:"avatar_url"`
	CreatedAt time.Time `json:"created_at"`
}

type LoginGoogleResponseData struct {
	AccessToken string       `json:"access_token"`
	User        UserResource `json:"user"`
}
