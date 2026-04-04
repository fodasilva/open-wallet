package transactions

import (
	"github.com/felipe1496/open-wallet/internal/constants"
	"github.com/felipe1496/open-wallet/internal/resources/transactions/repository"
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
	Transaction repository.Transaction `json:"transaction"`
}

type CreateTransactionResponse struct {
	Data CreateTransactionResponseData `json:"data"`
}

type CreateTransactionResponseData struct {
	Transaction repository.Transaction `json:"transaction"`
}

type ListEntriesResponse struct {
	Data  ListEntriesResponseData `json:"data"`
	Query utils.QueryMeta         `json:"query"`
}

type ListEntriesResponseData struct {
	Entries []repository.ViewEntry `json:"entries"`
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

// DELETED Transaction, Entry, ViewEntry models from here as they are now in the repository package.
