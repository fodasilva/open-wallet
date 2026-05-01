package usecases

import (
	"context"
	"net/http"

	"go.opentelemetry.io/otel"

	transactionRepo "github.com/felipe1496/open-wallet/internal/resources/transactions/repository"
	"github.com/felipe1496/open-wallet/internal/util"
)

func (uc *TransactionsUseCasesImpl) ListEntries(ctx context.Context) ([]transactionRepo.ViewEntry, error) {
	tracer := otel.Tracer("usecase")
	_, span := tracer.Start(ctx, "TransactionsUseCase.ListEntries")
	defer span.End()

	entries, err := uc.entriesRepo.Select(ctx, uc.db)

	if err != nil {
		span.RecordError(err)
		return []transactionRepo.ViewEntry{}, util.NewHTTPError(http.StatusInternalServerError, "failed to fetch entries")
	}

	return entries, nil
}
