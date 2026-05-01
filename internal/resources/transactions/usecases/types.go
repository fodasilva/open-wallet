package usecases

import (
	"time"

	transactionRepo "github.com/felipe1496/open-wallet/internal/resources/transactions/repository"
	"github.com/felipe1496/open-wallet/internal/util"
)

type CreateTransactionDTO struct {
	UserID       string
	Name         string
	CategoryID   util.OptionalNullable[string]
	Note         util.OptionalNullable[string]
	Type         transactionRepo.TransactionType
	Entries      []CreateEntryDTO
	RecurrenceID util.OptionalNullable[string]
}

type CreateEntryDTO struct {
	Amount        float64
	ReferenceDate time.Time
}

type UpdateTransactionDTO struct {
	Name         util.OptionalNullable[string]
	Note         util.OptionalNullable[string]
	CategoryID   util.OptionalNullable[string]
	Entries      util.OptionalNullable[[]UpdateEntryDTO]
	RecurrenceID util.OptionalNullable[string]
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
