package factory

import (
	"database/sql"

	"github.com/felipe1496/open-wallet/infra"
	authUseCases "github.com/felipe1496/open-wallet/internal/resources/auth/usecases"
	categoriesUseCases "github.com/felipe1496/open-wallet/internal/resources/categories/usecases"
	categoriesRepo "github.com/felipe1496/open-wallet/internal/resources/categories/repository"
	recurrencesUseCases "github.com/felipe1496/open-wallet/internal/resources/recurrences/usecases"
	recurrencesRepo "github.com/felipe1496/open-wallet/internal/resources/recurrences/repository"
	"github.com/felipe1496/open-wallet/internal/resources/transactions"
	transactionsRepo "github.com/felipe1496/open-wallet/internal/resources/transactions/repository"
	"github.com/felipe1496/open-wallet/internal/resources/users"
	usersRepo "github.com/felipe1496/open-wallet/internal/resources/users/repository"
	"github.com/felipe1496/open-wallet/internal/services"
)

type Factory struct {
	db  *sql.DB
	cfg *infra.Config

	googleService       services.GoogleService
	jwtService          services.JWTService
	usersUseCase        users.UsersUseCase
	authUseCases        authUseCases.AuthUseCases
	categoriesUseCases  categoriesUseCases.CategoriesUseCases
	transactionsUseCase transactions.TransactionsUseCase
	recurrencesUseCases  recurrencesUseCases.RecurrencesUseCases
}

func NewFactory(db *sql.DB, cfg *infra.Config) *Factory {
	return &Factory{db: db, cfg: cfg}
}

func (f *Factory) GoogleService() services.GoogleService {
	if f.googleService == nil {
		f.googleService = services.NewGoogleService(f.cfg)
	}
	return f.googleService
}

func (f *Factory) JWTService() services.JWTService {
	if f.jwtService == nil {
		f.jwtService = services.NewJWTService(f.cfg)
	}
	return f.jwtService
}

func (f *Factory) UsersUseCase() users.UsersUseCase {
	if f.usersUseCase == nil {
		f.usersUseCase = users.NewUsersUseCase(usersRepo.NewUsersRepo(), f.db)
	}
	return f.usersUseCase
}

func (f *Factory) AuthUseCases() authUseCases.AuthUseCases {
	if f.authUseCases == nil {
		f.authUseCases = authUseCases.NewAuthUseCases(f.GoogleService(), f.UsersUseCase())
	}
	return f.authUseCases
}

func (f *Factory) CategoriesUseCases() categoriesUseCases.CategoriesUseCases {
	if f.categoriesUseCases == nil {
		f.categoriesUseCases = categoriesUseCases.NewCategoriesUseCases(categoriesRepo.NewCategoriesRepo(), f.db)
	}
	return f.categoriesUseCases
}

func (f *Factory) TransactionsUseCase() transactions.TransactionsUseCase {
	if f.transactionsUseCase == nil {
		f.transactionsUseCase = transactions.NewTransactionsUseCase(
			transactionsRepo.NewTransactionsRepo(),
			transactionsRepo.NewEntriesRepo(),
			f.CategoriesUseCases(),
			f.db,
		)
	}
	return f.transactionsUseCase
}

func (f *Factory) RecurrencesUseCases() recurrencesUseCases.RecurrencesUseCases {
	if f.recurrencesUseCases == nil {
		f.recurrencesUseCases = recurrencesUseCases.NewRecurrencesUseCases(
			recurrencesRepo.NewRecurrencesRepo(),
			f.CategoriesUseCases(),
			f.TransactionsUseCase(),
			f.db,
		)
	}
	return f.recurrencesUseCases
}
