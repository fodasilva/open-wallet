package handlers

import (
	"time"

	transactionRepo "github.com/felipe1496/open-wallet/internal/resources/transactions/repository"
)

type CreateTransactionRequest struct {
	Name       string                          `json:"name" binding:"required,min=1,max=100"`
	CategoryID *string                         `json:"category_id" binding:"omitempty"`
	Note       *string                         `json:"note" binding:"omitempty,min=0,max=400"`
	Type       transactionRepo.TransactionType `json:"type" binding:"required,oneof=installment simple_expense income"`
	Entries    []CreateEntryRequest            `json:"entries" binding:"required,min=1,max=100,dive"`
}

type CreateEntryRequest struct {
	Amount        float64 `json:"amount" binding:"required,gte=-999999,lte=999999"`
	ReferenceDate string  `json:"reference_date" binding:"required,datetime=2006-01-02"`
}

type UpdateTransactionRequest struct {
	Name       *string               `json:"name,omitempty" binding:"omitempty,min=1,max=100"`
	CategoryID *string               `json:"category_id,omitempty" binding:"omitempty"`
	Note       *string               `json:"note,omitempty" binding:"omitempty,min=0,max=400"`
	Entries    *[]UpdateEntryRequest `json:"entries,omitempty" binding:"omitempty,min=1,max=100,dive"`
}

type UpdateEntryRequest struct {
	Amount        float64 `json:"amount" binding:"required,gte=-999999,lte=999999"`
	ReferenceDate string  `json:"reference_date" binding:"required,datetime=2006-01-02"`
}

type TransactionResource struct {
	ID            string    `json:"id" binding:"required"`
	UserID        string    `json:"user_id" binding:"required"`
	Type          string    `json:"type" binding:"required"`
	Name          string    `json:"name" binding:"required"`
	Description   *string   `json:"description"`
	CreatedAt     time.Time `json:"created_at" binding:"required"`
	CategoryID    *string   `json:"category_id"`
	CategoryName  *string   `json:"category_name,omitempty"`
	CategoryColor *string   `json:"category_color,omitempty"`
	RecurrenceID  *string   `json:"recurrence_id"`
}

type EntryResource struct {
	ID                string    `json:"id" binding:"required"`
	TransactionID     string    `json:"transaction_id" binding:"required"`
	Name              string    `json:"name" binding:"required"`
	Description       *string   `json:"description"`
	Amount            float64   `json:"amount" binding:"required"`
	Period            string    `json:"period" binding:"required"`
	UserID            string    `json:"user_id" binding:"required"`
	Type              string    `json:"type" binding:"required"`
	TotalAmount       float64   `json:"total_amount" binding:"required"`
	Installment       int       `json:"installment" binding:"required"`
	TotalInstallments int       `json:"total_installments" binding:"required"`
	CreatedAt         time.Time `json:"created_at" binding:"required"`
	ReferenceDate     time.Time `json:"reference_date" binding:"required"`
	CategoryID        *string   `json:"category_id,omitempty"`
	CategoryName      *string   `json:"category_name,omitempty"`
	CategoryColor     *string   `json:"category_color,omitempty"`
	RecurrenceID      *string   `json:"recurrence_id,omitempty"`
}

type UpdateTransactionResponseData struct {
	Transaction TransactionResource `json:"transaction" binding:"required"`
}

type CreateTransactionResponseData struct {
	Transaction TransactionResource `json:"transaction" binding:"required"`
}

type ListEntriesResponseData struct {
	Entries []EntryResource `json:"entries" binding:"required"`
}

type MonthlySummaryResource struct {
	Period  string  `json:"period" binding:"required"`
	Income  float64 `json:"income" binding:"required"`
	Expense float64 `json:"expense" binding:"required"`
	Balance float64 `json:"balance" binding:"required"`
}

type SummaryResponseData struct {
	Summary []MonthlySummaryResource `json:"summary" binding:"required"`
}
