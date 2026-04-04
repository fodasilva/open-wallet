package usecases

import (
	"context"
	"net/http"

	"github.com/felipe1496/open-wallet/internal/resources/recurrences/repository"
	"github.com/felipe1496/open-wallet/internal/utils"
	"go.opentelemetry.io/otel"
)

func (uc *RecurrencesUseCasesImpl) List(ctx context.Context, filter *utils.QueryOptsBuilder) ([]repository.Recurrence, error) {
	tracer := otel.Tracer("usecase")
	ctx, span := tracer.Start(ctx, "RecurrencesUseCase.List")
	defer span.End()

	items, err := uc.repo.Select(uc.db, filter)
	if err != nil {
		span.RecordError(err)
		return nil, utils.NewHTTPError(http.StatusInternalServerError, "failed to list recurrences")
	}

	return items, nil
}
