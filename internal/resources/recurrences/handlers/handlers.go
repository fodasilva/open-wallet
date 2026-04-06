package handlers

import (
	"github.com/felipe1496/open-wallet/internal/resources/recurrences/usecases"
)

type API struct {
	recurrencesUseCases usecases.RecurrencesUseCases
}

func NewHandler(recurrencesUseCases usecases.RecurrencesUseCases) *API {
	return &API{
		recurrencesUseCases: recurrencesUseCases,
	}
}
