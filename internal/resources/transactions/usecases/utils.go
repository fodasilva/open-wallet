package usecases

import (
	"net/http"
	"time"

	transactionRepo "github.com/felipe1496/open-wallet/internal/resources/transactions/repository"
	"github.com/felipe1496/open-wallet/internal/util/httputil"
)

type validateTransactionPropsEntry struct {
	Amount        float64
	ReferenceDate time.Time
}

func validateTransaction(entries []validateTransactionPropsEntry, transactionType transactionRepo.TransactionType) error {
	if err := validateEntriesCount(entries, transactionType); err != nil {
		return err
	}

	if err := validatePeriodsUniqueness(entries); err != nil {
		return err
	}

	return validateAmountsSigns(entries, transactionType)
}

func validateEntriesCount(entries []validateTransactionPropsEntry, tType transactionRepo.TransactionType) error {
	switch tType {
	case transactionRepo.SimpleExpense, transactionRepo.Income:
		if len(entries) > 1 {
			msg := "expense must have only one entry"
			if tType == transactionRepo.Income {
				msg = "income must have only one entry"
			}
			return httputil.NewHTTPError(http.StatusBadRequest, msg)
		}
	case transactionRepo.Installment:
		if len(entries) < 2 {
			return httputil.NewHTTPError(http.StatusBadRequest, "installment must have at least two entries")
		}
	}
	return nil
}

func validatePeriodsUniqueness(entries []validateTransactionPropsEntry) error {
	periods := make(map[string]bool)
	for _, entry := range entries {
		period := entry.ReferenceDate.Format("200601")
		if periods[period] {
			return httputil.NewHTTPError(http.StatusBadRequest, "entries must be in different periods")
		}
		periods[period] = true
	}
	return nil
}

func validateAmountsSigns(entries []validateTransactionPropsEntry, tType transactionRepo.TransactionType) error {
	for _, entry := range entries {
		switch tType {
		case transactionRepo.Installment, transactionRepo.SimpleExpense, transactionRepo.Recurrence:
			if entry.Amount >= 0 {
				msg := "installment entries must have amount lower than zero"
				if tType == transactionRepo.SimpleExpense {
					msg = "expense entry must have amount lower than zero"
				} else if tType == transactionRepo.Recurrence {
					msg = "recurrence entries must have amount lower than zero"
				}
				return httputil.NewHTTPError(http.StatusBadRequest, msg)
			}
		case transactionRepo.Income:
			if entry.Amount <= 0 {
				return httputil.NewHTTPError(http.StatusBadRequest, "income entry must have amount greater than zero")
			}
		}
	}
	return nil
}
