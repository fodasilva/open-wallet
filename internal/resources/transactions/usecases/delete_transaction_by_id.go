package usecases

import (
	"net/http"

	"github.com/felipe1496/open-wallet/internal/utils"
)

func (uc *TransactionsUseCasesImpl) DeleteTransactionById(id string) error {
	transactionExists, err := uc.transactionsRepo.Select(uc.db, utils.QueryOpts().And("id", "eq", id))

	if err != nil {
		return utils.NewHTTPError(http.StatusInternalServerError, "an error occurred while fetching transactions")
	}

	if len(transactionExists) == 0 {
		return utils.NewHTTPError(http.StatusNotFound, "transaction not found")
	}

	err = uc.transactionsRepo.Delete(uc.db, utils.QueryOpts().And("id", "eq", id))

	if err != nil {
		return utils.NewHTTPError(http.StatusInternalServerError, "it was not possible to delete the transaction")
	}

	return nil
}
