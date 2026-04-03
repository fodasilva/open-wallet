package repository

import (
	"github.com/felipe1496/open-wallet/internal/utils"
)

// Repository interface. Make sure to include methods
// that you defined with @method tags in models.go and any other methods you need.
type RecurrencesRepo interface {
	Select(db utils.Executer, filter *utils.QueryOptsBuilder) ([]Recurrence, error)
	Insert(db utils.Executer, data CreateRecurrenceDTO) error
	Update(db utils.Executer, data UpdateRecurrenceDTO, filter *utils.QueryOptsBuilder) error
	Delete(db utils.Executer, filter *utils.QueryOptsBuilder) error
	Count(db utils.Executer, filter *utils.QueryOptsBuilder) (int, error)
}

// Implementation struct. Name must match @name tag in models.go
type RecurrencesRepoImpl struct {
}

func NewRecurrencesRepo() RecurrencesRepo {
	return &RecurrencesRepoImpl{}
}
