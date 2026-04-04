package factory

import (
	"database/sql"

	"github.com/felipe1496/open-wallet/infra"
	"github.com/felipe1496/open-wallet/internal/resources/auth/usecases"
	"github.com/felipe1496/open-wallet/internal/resources/categories"
	categoriesRepo "github.com/felipe1496/open-wallet/internal/resources/categories/repository"
	"github.com/felipe1496/open-wallet/internal/resources/recurrences"
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
	authUseCases        usecases.AuthUseCases
	categoriesUseCase   categories.CategoriesUseCase
	transactionsUseCase transactions.TransactionsUseCase
	recurrencesUseCase  recurrences.RecurrencesUseCase
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

func (f *Factory) AuthUseCases() usecases.AuthUseCases {
	if f.authUseCases == nil {
		f.authUseCases = usecases.NewAuthUseCases(f.GoogleService(), f.UsersUseCase())
	}
	return f.authUseCases
}

func (f *Factory) CategoriesUseCase() categories.CategoriesUseCase {
	if f.categoriesUseCase == nil {
		f.categoriesUseCase = categories.NewCategoriesUseCase(categoriesRepo.NewCategoriesRepo(), f.db)
	}
	return f.categoriesUseCase
}

func (f *Factory) TransactionsUseCase() transactions.TransactionsUseCase {
	if f.transactionsUseCase == nil {
		f.transactionsUseCase = transactions.NewTransactionsUseCase(
			transactionsRepo.NewTransactionsRepo(),
			transactionsRepo.NewEntriesRepo(),
			f.CategoriesUseCase(),
			f.db,
		)
	}
	return f.transactionsUseCase
}

func (f *Factory) RecurrencesUseCase() recurrences.RecurrencesUseCase {
	if f.recurrencesUseCase == nil {
		f.recurrencesUseCase = recurrences.NewRecurrencesUseCase(
			recurrencesRepo.NewRecurrencesRepo(),
			f.CategoriesUseCase(),
			f.TransactionsUseCase(),
			f.db,
		)
	}
	return f.recurrencesUseCase
}
