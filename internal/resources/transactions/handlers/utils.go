package handlers

import (
	transactionRepo "github.com/felipe1496/open-wallet/internal/resources/transactions/repository"
)

func MapTransactionResource(t transactionRepo.Transaction) TransactionResource {
	return TransactionResource{
		ID:           t.ID,
		UserID:       t.UserID,
		Type:         string(t.Type),
		Name:         t.Name,
		Description:  t.Description,
		CreatedAt:    t.CreatedAt,
		CategoryID:   t.CategoryID,
		RecurrenceID: t.RecurrenceID,
	}
}

func MapEntryResource(e transactionRepo.ViewEntry) EntryResource {
	return EntryResource{
		ID:                e.ID,
		TransactionID:     e.TransactionID,
		Name:              e.Name,
		Description:       e.Description,
		Amount:            e.Amount,
		Period:            e.Period,
		UserID:            e.UserID,
		Type:              string(e.Type),
		TotalAmount:       e.TotalAmount,
		Installment:       e.Installment,
		TotalInstallments: e.TotalInstallments,
		CreatedAt:         e.CreatedAt,
		ReferenceDate:     e.ReferenceDate,
		CategoryID:        e.CategoryID,
		CategoryName:      e.CategoryName,
		CategoryColor:     e.CategoryColor,
		RecurrenceID:      e.RecurrenceID,
	}
}
