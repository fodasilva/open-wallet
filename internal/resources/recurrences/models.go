package recurrences

import (
	"time"

	"github.com/felipe1496/open-wallet/internal/utils"
)

// ==============================================================================
//  1. HTTP MODELS
// ==============================================================================

type CreateRecurrenceRequest struct {
	Name        string  `json:"name" binding:"required,min=1,max=100"`
	CategoryID  *string `json:"category_id" binding:"omitempty"`
	Note        *string `json:"note" binding:"omitempty,min=0,max=400"`
	Amount      float64 `json:"amount" binding:"required,lt=0"`
	DayOfMonth  int     `json:"day_of_month" binding:"required,min=1,max=31"`
	StartPeriod string  `json:"start_period" binding:"required,len=6"`
	EndPeriod   *string `json:"end_period" binding:"omitempty,len=6"`
}

type UpdateRecurrenceRequest struct {
	Update      []string `json:"update" binding:"required,min=1,dive,oneof=name category_id note amount day_of_month start_period end_period"`
	Name        *string  `json:"name" binding:"omitempty,min=1,max=100"`
	CategoryID  *string  `json:"category_id" binding:"omitempty"`
	Note        *string  `json:"note" binding:"omitempty,min=0,max=400"`
	Amount      *float64 `json:"amount" binding:"omitempty,lt=0"`
	DayOfMonth  *int     `json:"day_of_month" binding:"omitempty,min=1,max=31"`
	StartPeriod *string  `json:"start_period" binding:"omitempty,len=6"`
	EndPeriod   *string  `json:"end_period" binding:"omitempty,len=6"`
}

type CreateRecurrenceResponse struct {
	Data CreateRecurrenceResponseData `json:"data"`
}

type CreateRecurrenceResponseData struct {
	Recurrence Recurrence `json:"recurrence"`
}

type UpdateRecurrenceResponse struct {
	Data UpdateRecurrenceResponseData `json:"data"`
}

type UpdateRecurrenceResponseData struct {
	Recurrence Recurrence `json:"recurrence"`
}

type ListRecurrencesResponse struct {
	Data  ListRecurrencesResponseData `json:"data"`
	Query utils.QueryMeta             `json:"query"`
}

type ListRecurrencesResponseData struct {
	Recurrences []Recurrence `json:"recurrences"`
}

// ==============================================================================
// 2. DTO MODELS
// ==============================================================================

type CreateRecurrenceDTO struct {
	UserID      string
	Name        string
	CategoryID  *string
	Note        *string
	Amount      float64
	DayOfMonth  int
	StartPeriod string
	EndPeriod   *string
}

type UpdateRecurrenceDTO struct {
	Update      []string
	Name        *string
	CategoryID  *string
	Note        *string
	Amount      *float64
	DayOfMonth  *int
	StartPeriod *string
	EndPeriod   *string
}

// ==============================================================================
// 3. DATABASE
// ==============================================================================

// Recurrences table record
type Recurrence struct {
	ID            string    `json:"id"`
	UserID        string    `json:"user_id"`
	Name          string    `json:"name"`
	Note          *string   `json:"note"`
	Amount        float64   `json:"amount"`
	DayOfMonth    int       `json:"day_of_month"`
	CategoryID    *string   `json:"category_id"`
	CategoryName  *string   `json:"category_name"`
	CategoryColor *string   `json:"category_color"`
	StartPeriod   string    `json:"start_period"`
	EndPeriod     *string   `json:"end_period"`
	CreatedAt     time.Time `json:"created_at"`
}
