package handlers

import (
	"github.com/felipe1496/open-wallet/internal/resources/transactions/usecases"
)

type API struct {
	transactionsUseCases usecases.TransactionsUseCases
}

func NewHandler(transactionsUseCases usecases.TransactionsUseCases) *API {
	return &API{
		transactionsUseCases: transactionsUseCases,
	}
}
