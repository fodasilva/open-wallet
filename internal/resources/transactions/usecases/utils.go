package usecases

import (
	transactionRepo "github.com/felipe1496/open-wallet/internal/resources/transactions/repository"
	"net/http"
	"time"

	"github.com/felipe1496/open-wallet/internal/utils"
)

type validateTransactionPropsEntry struct {
	Amount        float64
	ReferenceDate string
}

func validateTransaction(entries []validateTransactionPropsEntry, transactionType transactionRepo.TransactionType) error {
	switch transactionType {
	case transactionRepo.SimpleExpense, transactionRepo.Income:
		if len(entries) > 1 {
			msg := "expense must have only one entry"
			if transactionType == transactionRepo.Income {
				msg = "income must have only one entry"
			}
			return utils.NewHTTPError(http.StatusBadRequest, msg)
		}
	case transactionRepo.Installment:
		if len(entries) < 2 {
			return utils.NewHTTPError(http.StatusBadRequest, "installment must have at least two entries")
		}
	}

	for i, refEntry := range entries {
		iRefDate, _ := time.Parse("2006-01-02", refEntry.ReferenceDate)
		iPeriod := iRefDate.Format("200601")
		for j, currEntry := range entries {
			if i != j {
				jRefDate, _ := time.Parse("2006-01-02", currEntry.ReferenceDate)
				jPeriod := jRefDate.Format("200601")
				if iPeriod == jPeriod {
					return utils.NewHTTPError(http.StatusBadRequest, "entries must be in different periods")
				}
			}
		}

		switch transactionType {
		case transactionRepo.Installment, transactionRepo.SimpleExpense, transactionRepo.Recurrence:
			if refEntry.Amount >= 0 {
				msg := "installment entries must have amount lower than zero"
				if transactionType == transactionRepo.SimpleExpense {
					msg = "expense entry must have amount lower than zero"
				} else if transactionType == transactionRepo.Recurrence {
					msg = "recurrence entries must have amount lower than zero"
				}
				return utils.NewHTTPError(http.StatusBadRequest, msg)
			}
		case transactionRepo.Income:
			if refEntry.Amount <= 0 {
				return utils.NewHTTPError(http.StatusBadRequest, "income entry must have amount greater than zero")
			}
		}
	}
	return nil
}
