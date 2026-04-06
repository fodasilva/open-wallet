package handlers

import (
	"github.com/felipe1496/open-wallet/internal/resources/recurrences/repository"
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
