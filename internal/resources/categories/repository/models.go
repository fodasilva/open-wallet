package repository

import (
	"time"

	"github.com/felipe1496/open-wallet/internal/utils"
)

// @gen_repo
// @table: categories
// @entity: Category
// @name: CategoriesRepoImpl
// @method: Select | fields: id:ID, user_id:UserID, name:Name, color:Color, created_at:CreatedAt
// @method: Insert | fields: id:ID, user_id:UserID, name:Name, color:Color | payload: CreateCategoryDTO
// @method: Update | fields: name:Name?, color:Color? | payload: UpdateCategoryDTO
// @method: Delete
// @method: Count

type Category struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Name      string    `json:"name"`
	Color     string    `json:"color"`
	CreatedAt time.Time `json:"created_at"`
}

type CategoryAmountPerPeriod struct {
	ID          string  `json:"id"`
	UserID      string  `json:"user_id"`
	Name        string  `json:"name"`
	Color       string  `json:"color"`
	Period      string  `json:"period"`
	TotalAmount float64 `json:"total_amount"`
}

type CreateCategoryDTO struct {
	ID     string
	UserID string
	Name   string
	Color  string
}

type UpdateCategoryDTO struct {
	Name  utils.OptionalNullable[string]
	Color utils.OptionalNullable[string]
}
