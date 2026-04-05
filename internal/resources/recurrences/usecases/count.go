package usecases

import (
	"context"
	"net/http"

	"go.opentelemetry.io/otel"

	"github.com/felipe1496/open-wallet/internal/utils"
	"github.com/felipe1496/open-wallet/internal/utils/querybuilder"
)

func (uc *RecurrencesUseCasesImpl) Count(ctx context.Context, filter *querybuilder.Builder) (int, error) {
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
