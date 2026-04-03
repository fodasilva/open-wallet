package transactions

import (
	"time"

	"github.com/felipe1496/open-wallet/internal/constants"
	"github.com/felipe1496/open-wallet/internal/utils"
)

// ==============================================================================
//  1. HTTP MODELS
//     Models that represents request or response objects
//
// ==============================================================================
type CreateTransactionRequest struct {
	Name       string                    `json:"name" binding:"required,min=1,max=100"`
	CategoryID *string                   `json:"category_id" binding:"omitempty"`
	Note       *string                   `json:"note" binding:"omitempty,min=0,max=400"`
	Type       constants.TransactionType `json:"type" binding:"required,oneof=installment simple_expense income"`
	Entries    []CreateEntryRequest      `json:"entries" binding:"required,min=1,max=100,dive"`
}

type CreateEntryRequest struct {
	Amount        float64 `json:"amount" binding:"required,gte=-999999,lte=999999"`
	ReferenceDate string  `json:"reference_date" binding:"required,datetime=2006-01-02"`
}

type UpdateTransactionRequest struct {
	Update     []string              `json:"update" binding:"required,min=1,dive,oneof=name category_id note entries"`
	Name       *string               `json:"name" binding:"omitempty,min=1,max=100"`
	CategoryID *string               `json:"category_id" binding:"omitempty"`
	Note       *string               `json:"note" binding:"omitempty,min=0,max=400"`
	Entries    *[]UpdateEntryRequest `json:"entries" binding:"omitempty,min=1,max=100,dive"`
}

type UpdateEntryRequest struct {
	Amount        float64 `json:"amount" binding:"required,gte=-999999,lte=999999"`
	ReferenceDate string  `json:"reference_date" binding:"required,datetime=2006-01-02"`
}

type UpdateTransactionResponse struct {
	Data UpdateTransactionResponseData `json:"data"`
}

type UpdateTransactionResponseData struct {
	Transaction Transaction `json:"transaction"`
}

type CreateTransactionResponse struct {
	Data CreateTransactionResponseData `json:"data"`
}

type CreateTransactionResponseData struct {
	Transaction Transaction `json:"transaction"`
}

type ListEntriesResponse struct {
	Data  ListEntriesResponseData `json:"data"`
	Query utils.QueryMeta         `json:"query"`
}

type ListEntriesResponseData struct {
	Entries []ViewEntry `json:"entries"`
}

// ==============================================================================
// 2. DTO MODELS
//    Models that represents data transfer objects between api layers
// ==============================================================================

type CreateTransactionDTO struct {
	UserID       string
	Name         string
	CategoryID   *string
	Note         *string
	Type         constants.TransactionType
	Entries      []CreateEntryDTO
	RecurrenceID *string
}

type CreateEntryDTO struct {
	Amount        float64
	ReferenceDate string
}

type UpdateTransactionDTO struct {
	Update       []string
	Name         *string
	Note         *string
	CategoryID   *string
	Entries      *[]UpdateEntryDTO
	RecurrenceID *string
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

// ==============================================================================
// 3. DATABASE
//    Models that represents database objects
// ==============================================================================

// View that mixes the entries with the transaction information, riched with some valuable information about the totality of this relationship
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

// Entries table record
type Entry struct {
	ID            string
	TransactionID string
	Amount        float64
	ReferenceDate string
	CreatedAt     time.Time
}

// Transactions table record
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
