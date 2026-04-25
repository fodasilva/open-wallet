package handlers

import (
	"time"
)

type CreateCategoryRequest struct {
	Name  string `json:"name" binding:"required"`
	Color string `json:"color" binding:"required"`
}

type CategoryResource struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Name      string    `json:"name"`
	Color     string    `json:"color"`
	CreatedAt time.Time `json:"created_at"`
}

type CategoryAmountPerPeriodResource struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Color       string  `json:"color"`
	TotalAmount float64 `json:"total_amount"`
}

type CreateCategoryResponseData struct {
	Category CategoryResource `json:"category"`
}

type ListCategoriesResponseData struct {
	Categories []CategoryResource `json:"categories"`
}

type UpdateCategoryRequest struct {
	Name  *string `json:"name" binding:"omitempty,min=1,max=50"`
	Color *string `json:"color" binding:"omitempty,hexcolor,len=7"`
}

type UpdateCategoryResponseData struct {
	Category CategoryResource `json:"category"`
}

type ListCategoryAmountPerPeriodResponseData struct {
	Categories []CategoryAmountPerPeriodResource `json:"categories"`
}
