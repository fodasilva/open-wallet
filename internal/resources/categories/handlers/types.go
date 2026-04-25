package handlers

import (
	"time"
)

type CreateCategoryRequest struct {
	Name  string `json:"name" binding:"required"`
	Color string `json:"color" binding:"required"`
}

type CategoryResource struct {
	ID        string    `json:"id" binding:"required"`
	UserID    string    `json:"user_id" binding:"required"`
	Name      string    `json:"name" binding:"required"`
	Color     string    `json:"color" binding:"required"`
	CreatedAt time.Time `json:"created_at" binding:"required"`
}

type CategoryAmountPerPeriodResource struct {
	ID          string  `json:"id" binding:"required"`
	Name        string  `json:"name" binding:"required"`
	Color       string  `json:"color" binding:"required"`
	TotalAmount float64 `json:"total_amount" binding:"required"`
}

type CreateCategoryResponseData struct {
	Category CategoryResource `json:"category" binding:"required"`
}

type ListCategoriesResponseData struct {
	Categories []CategoryResource `json:"categories" binding:"required"`
}

type UpdateCategoryRequest struct {
	Name  *string `json:"name" binding:"omitempty,min=1,max=50"`
	Color *string `json:"color" binding:"omitempty,hexcolor,len=7"`
}

type UpdateCategoryResponseData struct {
	Category CategoryResource `json:"category" binding:"required"`
}

type ListCategoryAmountPerPeriodResponseData struct {
	Categories []CategoryAmountPerPeriodResource `json:"categories" binding:"required"`
}
