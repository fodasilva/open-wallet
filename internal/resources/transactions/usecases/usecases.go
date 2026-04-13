package usecases

import (
	"context"
	"database/sql"

	categories "github.com/felipe1496/open-wallet/internal/resources/categories/usecases"
	transactionRepo "github.com/felipe1496/open-wallet/internal/resources/transactions/repository"
)

type TransactionsUseCases interface {
	ListEntries(ctx context.Context) ([]transactionRepo.ViewEntry, error)
	CountEntries(ctx context.Context) (int, error)
	DeleteTransactionById(ctx context.Context, id string, userID string) error
	CreateTransaction(ctx context.Context, payload CreateTransactionDTO) (transactionRepo.Transaction, error)
	UpdateTransaction(ctx context.Context, transactionID string, userID string, payload UpdateTransactionDTO) (transactionRepo.Transaction, error)
}

type TransactionsUseCasesImpl struct {
	transactionsRepo  transactionRepo.TransactionsRepo
	entriesRepo       transactionRepo.EntriesRepo
	categoriesUseCase categories.CategoriesUseCases
	db                *sql.DB
}

func NewTransactionsUseCases(
	transactionsRepo transactionRepo.TransactionsRepo,
	entriesRepo transactionRepo.EntriesRepo,
	categoriesUseCase categories.CategoriesUseCases,
	db *sql.DB,
) TransactionsUseCases {

	return &TransactionsUseCasesImpl{
		transactionsRepo:  transactionsRepo,
		entriesRepo:       entriesRepo,
		categoriesUseCase: categoriesUseCase,
		db:                db,
	}
}
