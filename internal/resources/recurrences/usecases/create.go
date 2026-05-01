package usecases

import (
	"context"
	"net/http"

	"github.com/oklog/ulid/v2"

	"github.com/felipe1496/open-wallet/internal/resources/recurrences/repository"
	"github.com/felipe1496/open-wallet/internal/util/httputil"
	"github.com/felipe1496/open-wallet/internal/util/querybuilder"
)

func (uc *RecurrencesUseCasesImpl) Create(ctx context.Context, payload repository.CreateRecurrenceDTO) (repository.Recurrence, error) {
	if payload.CategoryID.Set && payload.CategoryID.Value != nil {
		filterCtx := querybuilder.WithBuilder(ctx, querybuilder.New().
			And("id", "eq", *payload.CategoryID.Value).
			And("user_id", "eq", payload.UserID))
		categoryExists, err := uc.categoriesUseCase.List(filterCtx)
		if err != nil {
			return repository.Recurrence{}, httputil.NewHTTPError(http.StatusInternalServerError, "failed to check if category exists")
		}

		if len(categoryExists) == 0 {
			return repository.Recurrence{}, httputil.NewHTTPError(http.StatusNotFound, "category not found")
		}
	}

	if payload.ID == "" {
		payload.ID = ulid.Make().String()
	}

	err := uc.repo.Insert(ctx, uc.db, payload)
	if err != nil {
		return repository.Recurrence{}, httputil.NewHTTPError(http.StatusInternalServerError, "failed to create recurrence")
	}

	createdCtx := querybuilder.WithBuilder(ctx, querybuilder.New().And("id", "eq", payload.ID))
	recs, err := uc.repo.Select(createdCtx, uc.db)
	if err != nil || len(recs) == 0 {
		return repository.Recurrence{}, httputil.NewHTTPError(http.StatusInternalServerError, "failed to fetch created recurrence")
	}

	return recs[0], nil
}
