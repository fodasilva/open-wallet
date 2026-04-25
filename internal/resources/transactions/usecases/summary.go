package usecases

import (
	"context"

	transactionRepo "github.com/felipe1496/open-wallet/internal/resources/transactions/repository"
)

func (u *TransactionsUseCasesImpl) Summary(ctx context.Context) ([]transactionRepo.ViewSummary, error) {
	return u.summariesRepo.Select(ctx, u.db)
}
