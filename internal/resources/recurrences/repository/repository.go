package repository

import (
	"context"

	"github.com/felipe1496/open-wallet/internal/util"
)

// Repository interface. Make sure to include methods
// that you defined with @method tags in types.go and any other methods you need.
type RecurrencesRepo interface {
	Select(ctx context.Context, db util.Executer) ([]Recurrence, error)
	Insert(ctx context.Context, db util.Executer, data CreateRecurrenceDTO) error
	Update(ctx context.Context, db util.Executer, data UpdateRecurrenceDTO) error
	Delete(ctx context.Context, db util.Executer) error
	Count(ctx context.Context, db util.Executer) (int, error)
}

// Implementation struct. Name must match @name tag in types.go
type RecurrencesRepoImpl struct {
}

func NewRecurrencesRepo() RecurrencesRepo {
	return &RecurrencesRepoImpl{}
}
