package handlers

import (
	"github.com/felipe1496/open-wallet/internal/resources/recurrences/repository"
	"github.com/felipe1496/open-wallet/internal/utils"
)

type CreateRecurrenceRequest struct {
	Name        string  `json:"name" binding:"required,min=1,max=100"`
	CategoryID  *string `json:"category_id" binding:"omitempty"`
	Note        *string `json:"note" binding:"omitempty,min=0,max=400"`
	Amount      float64 `json:"amount" binding:"required,lt=0"`
	DayOfMonth  int     `json:"day_of_month" binding:"required,min=1,max=31"`
	StartPeriod string  `json:"start_period" binding:"required,len=6"`
	EndPeriod   *string `json:"end_period" binding:"omitempty,len=6"`
}

type CreateRecurrenceResponse struct {
	Data CreateRecurrenceResponseData `json:"data"`
}

type CreateRecurrenceResponseData struct {
	Recurrence repository.Recurrence `json:"recurrence"`
}

type UpdateRecurrenceRequest struct {
	Name        *string  `json:"name" binding:"omitempty,min=1,max=100"`
	CategoryID  *string  `json:"category_id" binding:"omitempty"`
	Note        *string  `json:"note" binding:"omitempty,min=0,max=400"`
	Amount      *float64 `json:"amount" binding:"omitempty,lt=0"`
	DayOfMonth  *int     `json:"day_of_month" binding:"omitempty,min=1,max=31"`
	StartPeriod *string  `json:"start_period" binding:"omitempty,len=6"`
	EndPeriod   *string  `json:"end_period" binding:"omitempty,len=6"`
}

type UpdateRecurrenceResponse struct {
	Data UpdateRecurrenceResponseData `json:"data"`
}

type UpdateRecurrenceResponseData struct {
	Recurrence repository.Recurrence `json:"recurrence"`
}

type ListRecurrencesResponse struct {
	Data  ListRecurrencesResponseData `json:"data"`
	Query utils.QueryMeta             `json:"query"`
}

type ListRecurrencesResponseData struct {
	Recurrences []repository.Recurrence `json:"recurrences"`
}
