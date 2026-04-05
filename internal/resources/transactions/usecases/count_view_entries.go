package usecases

import (
	"context"
	"net/http"

	"go.opentelemetry.io/otel"

	"github.com/felipe1496/open-wallet/internal/utils"
)

func (uc *TransactionsUseCasesImpl) CountEntries(ctx context.Context, filter *utils.QueryOptsBuilder) (int, error) {
	tracer := otel.Tracer("usecase")
	_, span := tracer.Start(ctx, "TransactionsUseCase.CountEntries")
	defer span.End()

	count, err := uc.entriesRepo.Count(uc.db, filter)

	if err != nil {
		span.RecordError(err)
		return 0, utils.NewHTTPError(http.StatusInternalServerError, "failed to count entries")
	}

	return count, nil
}
