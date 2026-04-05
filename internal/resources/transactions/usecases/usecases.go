package usecases

import (
	"context"
	"database/sql"

	categories "github.com/felipe1496/open-wallet/internal/resources/categories/usecases"
	transactionRepo "github.com/felipe1496/open-wallet/internal/resources/transactions/repository"
	"github.com/felipe1496/open-wallet/internal/utils"
)

type TransactionsUseCases interface {
	ListEntries(ctx context.Context, filter *utils.QueryOptsBuilder) ([]transactionRepo.ViewEntry, error)
	CountEntries(ctx context.Context, filter *utils.QueryOptsBuilder) (int, error)
	DeleteTransactionById(id string) error
	CreateTransaction(payload CreateTransactionDTO) (transactionRepo.Transaction, error)
	UpdateTransaction(transactionID string, userID string, payload UpdateTransactionDTO) (transactionRepo.Transaction, error)
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
