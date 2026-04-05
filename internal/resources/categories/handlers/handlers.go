package handlers

import (
	"github.com/felipe1496/open-wallet/internal/resources/categories/usecases"
)

type API struct {
	categoriesUseCases usecases.CategoriesUseCases
}

func NewHandler(categoriesUseCases usecases.CategoriesUseCases) *API {
	return &API{
		categoriesUseCases: categoriesUseCases,
	}
}
