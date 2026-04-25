package handlers

import (
	"github.com/felipe1496/open-wallet/internal/resources/categories/repository"
)

func MapCategoryResource(c repository.Category) CategoryResource {
	return CategoryResource{
		ID:        c.ID,
		UserID:    c.UserID,
		Name:      c.Name,
		Color:     c.Color,
		CreatedAt: c.CreatedAt,
	}
}

func MapCategoryAmountPerPeriodResource(c repository.CategoryAmountPerPeriod) CategoryAmountPerPeriodResource {
	return CategoryAmountPerPeriodResource{
		ID:          c.ID,
		Name:        c.Name,
		Color:       c.Color,
		TotalAmount: c.TotalAmount,
	}
}
