package repository

import (
	"time"

	"github.com/felipe1496/open-wallet/internal/util"
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
	ID           string
	UserID       string
	Type         TransactionType
	Name         string
	Description  *string
	CreatedAt    time.Time
	CategoryID   *string
	RecurrenceID *string
}

type CreateTransactionDTO struct {
	ID           string
	UserID       string
	Name         string
	CategoryID   util.OptionalNullable[string]
	Note         util.OptionalNullable[string]
	Type         TransactionType
	RecurrenceID util.OptionalNullable[string]
}

type UpdateTransactionDTO struct {
	Name         util.OptionalNullable[string]
	Note         util.OptionalNullable[string]
	CategoryID   util.OptionalNullable[string]
	RecurrenceID util.OptionalNullable[string]
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
	ReferenceDate time.Time
	CreatedAt     time.Time
}

type CreateEntryDTO struct {
	ID            string
	TransactionID string
	Amount        float64
	ReferenceDate time.Time
}

type UpdateEntryDTO struct {
	Amount        util.OptionalNullable[float64]
	ReferenceDate util.OptionalNullable[time.Time]
}

// @gen_repo
// @table: v_entries
// @entity: ViewEntry
// @name: EntriesRepoImpl
// @method: Select | fields: id:ID, transaction_id:TransactionID, name:Name, description:Description, amount:Amount, period:Period, user_id:UserID, category:Type, total_amount:TotalAmount, installment:Installment, total_installments:TotalInstallments, created_at:CreatedAt, reference_date:ReferenceDate, category_id:CategoryID, category_name:CategoryName, category_color:CategoryColor, recurrence_id:RecurrenceID
// @method: Count

type ViewEntry struct {
	ID                string
	TransactionID     string
	Name              string
	Description       *string
	Amount            float64
	Period            string
	UserID            string
	Type              TransactionType
	TotalAmount       float64
	Installment       int
	TotalInstallments int
	CreatedAt         time.Time
	ReferenceDate     time.Time
	CategoryID        *string
	CategoryName      *string
	CategoryColor     *string
	RecurrenceID      *string
}

// @gen_repo
// @table: v_summaries
// @entity: ViewSummary
// @name: SummariesRepoImpl
// @method: Select | fields: user_id:UserID, period:Period, total_expense:TotalExpense, total_income:TotalIncome, total_balance:TotalBalance

type ViewSummary struct {
	UserID       string  `json:"user_id"`
	Period       string  `json:"period"`
	TotalExpense float64 `json:"total_expense"`
	TotalIncome  float64 `json:"total_income"`
	TotalBalance float64 `json:"total_balance"`
}

type TransactionType string

const (
	SimpleExpense TransactionType = "simple_expense"
	Income        TransactionType = "income"
	Installment   TransactionType = "installment"
	Recurrence    TransactionType = "recurrence"
)
