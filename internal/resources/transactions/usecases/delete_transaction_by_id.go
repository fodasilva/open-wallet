package usecases

import (
	"context"
	"net/http"

	"github.com/felipe1496/open-wallet/internal/util/httputil"
	"github.com/felipe1496/open-wallet/internal/util/querybuilder"
)

func (uc *TransactionsUseCasesImpl) DeleteTransactionById(ctx context.Context, id string, userID string) error {
	filterCtx := querybuilder.WithBuilder(ctx, querybuilder.New().And("id", "eq", id).And("user_id", "eq", userID))
	transactionExists, err := uc.transactionsRepo.Select(filterCtx, uc.db)

	if err != nil {
		return httputil.NewHTTPError(http.StatusInternalServerError, "an error occurred while fetching transactions")
	}

	if len(transactionExists) == 0 {
		return httputil.NewHTTPError(http.StatusNotFound, "transaction not found")
	}

	err = uc.transactionsRepo.Delete(filterCtx, uc.db)

	if err != nil {
		return httputil.NewHTTPError(http.StatusInternalServerError, "it was not possible to delete the transaction")
	}

	return nil
}
