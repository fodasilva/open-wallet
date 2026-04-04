package repository

import (
	"github.com/felipe1496/open-wallet/internal/utils"
)

// Repository interface. Make sure to include methods
// that you defined with @method tags in types.go and any other methods you need.
type TransactionsRepo interface {
	Select(db utils.Executer, filter *utils.QueryOptsBuilder) ([]Transaction, error)
	Insert(db utils.Executer, data CreateTransactionDTO) error
	Update(db utils.Executer, data UpdateTransactionDTO, filter *utils.QueryOptsBuilder) error
	Delete(db utils.Executer, filter *utils.QueryOptsBuilder) error
}

// Implementation struct. Name must match @name tag in types.go
type TransactionsRepoImpl struct {
}

func NewTransactionsRepo() TransactionsRepo {
	return &TransactionsRepoImpl{}
}

// Repository interface. Make sure to include methods
// that you defined with @method tags in types.go and any other methods you need.
type EntriesRepo interface {
	Select(db utils.Executer, filter *utils.QueryOptsBuilder) ([]ViewEntry, error)
	Insert(db utils.Executer, data CreateEntryDTO) error
	Delete(db utils.Executer, filter *utils.QueryOptsBuilder) error
	Count(db utils.Executer, filter *utils.QueryOptsBuilder) (int, error)
}

// Implementation struct. Name must match @name tag in types.go
type EntriesRepoImpl struct {
}

func NewEntriesRepo() EntriesRepo {
	return &EntriesRepoImpl{}
}
