package usecases

import (
	"context"
	"net/http"

	"go.opentelemetry.io/otel"

	"github.com/felipe1496/open-wallet/internal/util"
)

func (uc *RecurrencesUseCasesImpl) Count(ctx context.Context) (int, error) {
	tracer := otel.Tracer("usecase")
	_, span := tracer.Start(ctx, "RecurrencesUseCase.Count")
	defer span.End()

	count, err := uc.repo.Count(ctx, uc.db)
	if err != nil {
		span.RecordError(err)
		return 0, util.NewHTTPError(http.StatusInternalServerError, "failed to count recurrences")
	}

	return count, nil
}
