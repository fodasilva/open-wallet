package transactions

import (
	transactionRepo "github.com/felipe1496/open-wallet/internal/resources/transactions/repository"
	"github.com/felipe1496/open-wallet/internal/utils"
)

// ==============================================================================
//  1. HTTP MODELS
//     Models that represents request or response objects
//
// ==============================================================================
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

type UpdateTransactionResponse struct {
	Data UpdateTransactionResponseData `json:"data"`
}

type UpdateTransactionResponseData struct {
	Transaction transactionRepo.Transaction `json:"transaction"`
}

type CreateTransactionResponse struct {
	Data CreateTransactionResponseData `json:"data"`
}

type CreateTransactionResponseData struct {
	Transaction transactionRepo.Transaction `json:"transaction"`
}

type ListEntriesResponse struct {
	Data  ListEntriesResponseData `json:"data"`
	Query utils.QueryMeta         `json:"query"`
}

type ListEntriesResponseData struct {
	Entries []transactionRepo.ViewEntry `json:"entries"`
}

// ==============================================================================
// 2. DTO MODELS
//    Models that represents data transfer objects between api layers
// ==============================================================================

// DELETED Transaction, Entry, ViewEntry models from here as they are now in the repository package.
