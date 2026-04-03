package repository

import (
	"time"

	"github.com/felipe1496/open-wallet/internal/utils"
)

// @gen_repo
// @table: recurrences
// @entity: Recurrence
// @name: RecurrencesRepoImpl
// @method: Insert | fields: id:ID, user_id:UserID, name:Name, note:Note?, amount:Amount, day_of_month:DayOfMonth, category_id:CategoryID?, start_period:StartPeriod, end_period:EndPeriod? | payload: CreateRecurrenceDTO
// @method: Update | fields: name:Name?, note:Note?, amount:Amount?, day_of_month:DayOfMonth?, category_id:CategoryID?, start_period:StartPeriod?, end_period:EndPeriod? | payload: UpdateRecurrenceDTO
// @method: Delete

// @gen_repo
// @table: v_recurrences
// @entity: Recurrence
// @name: RecurrencesRepoImpl
// @method: Select | fields: id:ID, user_id:UserID, name:Name, note:Note, amount:Amount, day_of_month:DayOfMonth, category_id:CategoryID, category_name:CategoryName, category_color:CategoryColor, start_period:StartPeriod, end_period:EndPeriod, created_at:CreatedAt
// @method: Count

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

type CreateRecurrenceDTO struct {
	ID          string
	UserID      string
	Name        string
	CategoryID  utils.OptionalNullable[string]
	Note        utils.OptionalNullable[string]
	Amount      float64
	DayOfMonth  int
	StartPeriod string
	EndPeriod   utils.OptionalNullable[string]
}

type UpdateRecurrenceDTO struct {
	Name        utils.OptionalNullable[string]
	CategoryID  utils.OptionalNullable[string]
	Note        utils.OptionalNullable[string]
	Amount      utils.OptionalNullable[float64]
	DayOfMonth  utils.OptionalNullable[int]
	StartPeriod utils.OptionalNullable[string]
	EndPeriod   utils.OptionalNullable[string]
}
