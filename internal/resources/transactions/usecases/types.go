package usecases

import (
	"github.com/felipe1496/open-wallet/internal/constants"
	"github.com/felipe1496/open-wallet/internal/utils"
)

type CreateTransactionDTO struct {
	UserID       string
	Name         string
	CategoryID   utils.OptionalNullable[string]
	Note         utils.OptionalNullable[string]
	Type         constants.TransactionType
	Entries      []CreateEntryDTO
	RecurrenceID utils.OptionalNullable[string]
}

type CreateEntryDTO struct {
	Amount        float64
	ReferenceDate string
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
	ReferenceDate string
}

type PersistEntryDTO struct {
	TransactionID string
	Amount        float64
	ReferenceDate string
}
