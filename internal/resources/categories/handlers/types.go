package handlers

import (
	"github.com/felipe1496/open-wallet/internal/resources/categories/repository"
	"github.com/felipe1496/open-wallet/internal/utils"
)

type CreateCategoryRequest struct {
	Name  string `json:"name" binding:"required"`
	Color string `json:"color" binding:"required"`
}

type CreateCategoryResponse struct {
	Data CreateCategoryResponseData `json:"data"`
}

type CreateCategoryResponseData struct {
	Category repository.Category `json:"category"`
}

type ListCategoriesResponse struct {
	Data  ListCategoriesResponseData `json:"data"`
	Query utils.QueryMeta            `json:"query"`
}

type ListCategoriesResponseData struct {
	Categories []repository.Category `json:"categories"`
}

type UpdateCategoryRequest struct {
	Name  *string `json:"name" binding:"omitempty,min=1,max=50"`
	Color *string `json:"color" binding:"omitempty,hexcolor,len=7"`
}

type UpdateCategoryResponse struct {
	Data UpdateCategoryResponseData `json:"data"`
}

type UpdateCategoryResponseData struct {
	Category repository.Category `json:"category"`
}

type ListCategoryAmountPerPeriodResponse struct {
	Data  ListCategoryAmountPerPeriodResponseData `json:"data"`
	Query utils.QueryMeta                         `json:"query"`
}

type ListCategoryAmountPerPeriodResponseData struct {
	Categories []repository.CategoryAmountPerPeriod `json:"categories"`
}

