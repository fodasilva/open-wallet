package repository

import (
	"context"

	"github.com/felipe1496/open-wallet/internal/utils"
)

// Repository interface. Make sure to include methods
// that you defined with @method tags in types.go and any other methods you need.
type TransactionsRepo interface {
	Select(ctx context.Context, db utils.Executer) ([]Transaction, error)
	Insert(ctx context.Context, db utils.Executer, data CreateTransactionDTO) error
	Update(ctx context.Context, db utils.Executer, data UpdateTransactionDTO) error
	Delete(ctx context.Context, db utils.Executer) error
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
	Select(ctx context.Context, db utils.Executer) ([]ViewEntry, error)
	Insert(ctx context.Context, db utils.Executer, data CreateEntryDTO) error
	Delete(ctx context.Context, db utils.Executer) error
	Count(ctx context.Context, db utils.Executer) (int, error)
}

// Implementation struct. Name must match @name tag in types.go
type EntriesRepoImpl struct {
}

func NewEntriesRepo() EntriesRepo {
	return &EntriesRepoImpl{}
}

type SummariesRepo interface {
	Select(ctx context.Context, db utils.Executer) ([]ViewSummary, error)
}

type SummariesRepoImpl struct {
}

func NewSummariesRepo() SummariesRepo {
	return &SummariesRepoImpl{}
}
