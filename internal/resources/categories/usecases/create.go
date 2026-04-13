package usecases

import (
	"context"
	"net/http"

	"github.com/oklog/ulid/v2"

	"github.com/felipe1496/open-wallet/internal/resources/categories/repository"
	"github.com/felipe1496/open-wallet/internal/utils"
)

func (uc *CategoriesUseCasesImpl) Create(ctx context.Context, payload repository.CreateCategoryDTO) (repository.Category, error) {
	if payload.ID == "" {
		payload.ID = ulid.Make().String()
	}

	err := uc.repo.Insert(ctx, uc.db, payload)

	if err != nil {
		return repository.Category{}, utils.NewHTTPError(http.StatusInternalServerError, "failed to create category")
	}

	return repository.Category{
		ID:     payload.ID,
		UserID: payload.UserID,
		Name:   payload.Name,
		Color:  payload.Color,
	}, nil
}
