package usecases

import (
	"time"

	transactionRepo "github.com/felipe1496/open-wallet/internal/resources/transactions/repository"
	"github.com/felipe1496/open-wallet/internal/utils"
)

type CreateTransactionDTO struct {
	UserID       string
	Name         string
	CategoryID   utils.OptionalNullable[string]
	Note         utils.OptionalNullable[string]
	Type         transactionRepo.TransactionType
	Entries      []CreateEntryDTO
	RecurrenceID utils.OptionalNullable[string]
}

type CreateEntryDTO struct {
	Amount        float64
	ReferenceDate time.Time
}

type UpdateTransactionDTO struct {
	Name         utils.OptionalNullable[string]
	Note         utils.OptionalNullable[string]
	CategoryID   utils.OptionalNullable[string]
	Entries      utils.OptionalNullable[[]UpdateEntryDTO]
	RecurrenceID utils.OptionalNullable[string]
}

type UpdateEntryDTO struct {
	Amount        float64
	ReferenceDate time.Time
}

type PersistEntryDTO struct {
	TransactionID string
	Amount        float64
	ReferenceDate time.Time
}
