package usecases

import (
	"net/http"

	"github.com/felipe1496/open-wallet/internal/resources/recurrences/repository"
	"github.com/felipe1496/open-wallet/internal/utils"
	"github.com/oklog/ulid/v2"
)

func (uc *RecurrencesUseCasesImpl) Create(payload repository.CreateRecurrenceDTO) (repository.Recurrence, error) {
	if payload.CategoryID.Set && payload.CategoryID.Value != nil {
		categoryExists, err := uc.categoriesUseCase.List(utils.QueryOpts().
			And("id", "eq", *payload.CategoryID.Value).
			And("user_id", "eq", payload.UserID))
		if err != nil {
			return repository.Recurrence{}, utils.NewHTTPError(http.StatusInternalServerError, "failed to check if category exists")
		}

		if len(categoryExists) == 0 {
			return repository.Recurrence{}, utils.NewHTTPError(http.StatusNotFound, "category not found")
		}
	}

	if payload.ID == "" {
		payload.ID = ulid.Make().String()
	}

	err := uc.repo.Insert(uc.db, payload)
	if err != nil {
		return repository.Recurrence{}, utils.NewHTTPError(http.StatusInternalServerError, "failed to create recurrence")
	}

	recs, err := uc.repo.Select(uc.db, utils.QueryOpts().And("id", "eq", payload.ID))
	if err != nil || len(recs) == 0 {
		return repository.Recurrence{}, utils.NewHTTPError(http.StatusInternalServerError, "failed to fetch created recurrence")
	}

	return recs[0], nil
}
