package repository

import (
	"time"

	"github.com/felipe1496/open-wallet/internal/constants"
	"github.com/felipe1496/open-wallet/internal/utils"
)

// @gen_repo
// @table: transactions
// @entity: Transaction
// @name: TransactionsRepoImpl
// @method: Select | fields: id:ID, user_id:UserID, category:Type, name:Name, description:Description, created_at:CreatedAt, category_id:CategoryID, recurrence_id:RecurrenceID
// @method: Insert | fields: id:ID, user_id:UserID, category:Type, name:Name, description:Note?, category_id:CategoryID?, recurrence_id:RecurrenceID? | payload: CreateTransactionDTO
// @method: Update | fields: name:Name?, description:Note?, category_id:CategoryID?, recurrence_id:RecurrenceID? | payload: UpdateTransactionDTO
// @method: Delete

type Transaction struct {
	ID           string                    `json:"id"`
	UserID       string                    `json:"user_id"`
	Type         constants.TransactionType `json:"type"`
	Name         string                    `json:"name"`
	Description  *string                   `json:"description"`
	CreatedAt    time.Time                 `json:"created_at"`
	CategoryID   *string                   `json:"category_id"`
	RecurrenceID *string                   `json:"recurrence_id"`
}

type CreateTransactionDTO struct {
	ID           string
	UserID       string
	Name         string
	CategoryID   utils.OptionalNullable[string]
	Note         utils.OptionalNullable[string]
	Type         constants.TransactionType
	RecurrenceID utils.OptionalNullable[string]
}

type UpdateTransactionDTO struct {
	Name         utils.OptionalNullable[string]
	Note         utils.OptionalNullable[string]
	CategoryID   utils.OptionalNullable[string]
	RecurrenceID utils.OptionalNullable[string]
}

// @gen_repo
// @table: entries
// @entity: Entry
// @name: EntriesRepoImpl
// @method: Insert | fields: id:ID, transaction_id:TransactionID, amount:Amount, reference_date:ReferenceDate | payload: CreateEntryDTO
// @method: Delete

type Entry struct {
	ID            string
	TransactionID string
	Amount        float64
	ReferenceDate string
	CreatedAt     time.Time
}

type CreateEntryDTO struct {
	ID            string
	TransactionID string
	Amount        float64
	ReferenceDate string
}

type UpdateEntryDTO struct {
	Amount        utils.OptionalNullable[float64]
	ReferenceDate utils.OptionalNullable[string]
}

// @gen_repo
// @table: v_entries
// @entity: ViewEntry
// @name: EntriesRepoImpl
// @method: Select | fields: id:ID, transaction_id:TransactionID, name:Name, description:Description, amount:Amount, period:Period, user_id:UserID, category:Type, total_amount:TotalAmount, installment:Installment, total_installments:TotalInstallments, created_at:CreatedAt, reference_date:ReferenceDate, category_id:CategoryID, category_name:CategoryName, category_color:CategoryColor, recurrence_id:RecurrenceID
// @method: Count

type ViewEntry struct {
	ID                string                    `json:"id"`
	TransactionID     string                    `json:"transaction_id"`
	Name              string                    `json:"name"`
	Description       *string                   `json:"description"`
	Amount            float64                   `json:"amount"`
	Period            string                    `json:"period"`
	UserID            string                    `json:"user_id"`
	Type              constants.TransactionType `json:"type"`
	TotalAmount       float64                   `json:"total_amount"`
	Installment       int                       `json:"installment"`
	TotalInstallments int                       `json:"total_installments"`
	CreatedAt         time.Time                 `json:"created_at"`
	ReferenceDate     string                    `json:"reference_date"`
	CategoryID        *string                   `json:"category_id,omitempty"`
	CategoryName      *string                   `json:"category_name,omitempty"`
	CategoryColor     *string                   `json:"category_color,omitempty"`
	RecurrenceID      *string                   `json:"recurrence_id,omitempty"`
}
