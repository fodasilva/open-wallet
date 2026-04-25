package handlers

import (
	"time"
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

type RecurrenceResource struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	Name        string    `json:"name"`
	Amount      float64   `json:"amount"`
	DayOfMonth  int       `json:"day_of_month"`
	StartPeriod string    `json:"start_period"`
	EndPeriod   *string   `json:"end_period"`
	Note        *string   `json:"note"`
	CategoryID  *string   `json:"category_id"`
	CreatedAt   time.Time `json:"created_at"`
}

type CreateRecurrenceResponseData struct {
	Recurrence RecurrenceResource `json:"recurrence"`
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

type UpdateRecurrenceResponseData struct {
	Recurrence RecurrenceResource `json:"recurrence"`
}

type ListRecurrencesResponseData struct {
	Recurrences []RecurrenceResource `json:"recurrences"`
}
