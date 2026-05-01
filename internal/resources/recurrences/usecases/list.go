package usecases

import (
	"context"
	"net/http"

	"go.opentelemetry.io/otel"

	"github.com/felipe1496/open-wallet/internal/resources/recurrences/repository"
	"github.com/felipe1496/open-wallet/internal/util/httputil"
)

func (uc *RecurrencesUseCasesImpl) List(ctx context.Context) ([]repository.Recurrence, error) {
	tracer := otel.Tracer("usecase")
	_, span := tracer.Start(ctx, "RecurrencesUseCase.List")
	defer span.End()

	items, err := uc.repo.Select(ctx, uc.db)
	if err != nil {
		span.RecordError(err)
		return nil, httputil.NewHTTPError(http.StatusInternalServerError, "failed to list recurrences")
	}

	return items, nil
}
