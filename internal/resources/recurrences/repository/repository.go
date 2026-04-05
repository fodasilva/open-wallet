package repository

import (
	"github.com/felipe1496/open-wallet/internal/utils"
	"github.com/felipe1496/open-wallet/internal/utils/querybuilder"
)

// Repository interface. Make sure to include methods
// that you defined with @method tags in types.go and any other methods you need.
type RecurrencesRepo interface {
	Select(db utils.Executer, filter *querybuilder.Builder) ([]Recurrence, error)
	Insert(db utils.Executer, data CreateRecurrenceDTO) error
	Update(db utils.Executer, data UpdateRecurrenceDTO, filter *querybuilder.Builder) error
	Delete(db utils.Executer, filter *querybuilder.Builder) error
	Count(db utils.Executer, filter *querybuilder.Builder) (int, error)
}

// Implementation struct. Name must match @name tag in types.go
type RecurrencesRepoImpl struct {
}

func NewRecurrencesRepo() RecurrencesRepo {
	return &RecurrencesRepoImpl{}
}
