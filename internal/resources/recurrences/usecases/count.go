package usecases

import (
	"context"
	"net/http"

	"github.com/felipe1496/open-wallet/internal/utils"
	"go.opentelemetry.io/otel"
)

func (uc *RecurrencesUseCasesImpl) Count(ctx context.Context, filter *utils.QueryOptsBuilder) (int, error) {
	tracer := otel.Tracer("usecase")
	_, span := tracer.Start(ctx, "RecurrencesUseCase.Count")
	defer span.End()

	count, err := uc.repo.Count(uc.db, filter)
	if err != nil {
		span.RecordError(err)
		return 0, utils.NewHTTPError(http.StatusInternalServerError, "failed to count recurrences")
	}

	return count, nil
}
