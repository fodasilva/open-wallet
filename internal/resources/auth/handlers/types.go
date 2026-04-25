package handlers

import (
	"time"
)

type LoginGoogleRequest struct {
	Code string `json:"code" binding:"required"`
}

type UserResource struct {
	ID        string    `json:"id" binding:"required"`
	Email     string    `json:"email" binding:"required"`
	Name      string    `json:"name" binding:"required"`
	AvatarURL string    `json:"avatar_url" binding:"required"`
	CreatedAt time.Time `json:"created_at" binding:"required"`
}

type LoginGoogleResponseData struct {
	AccessToken string       `json:"access_token" binding:"required"`
	User        UserResource `json:"user" binding:"required"`
}
