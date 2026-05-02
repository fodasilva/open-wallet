package factory

import (
	"database/sql"

	"github.com/felipe1496/open-wallet/infra"
	authUseCases "github.com/felipe1496/open-wallet/internal/resources/auth/usecases"
	categoriesRepo "github.com/felipe1496/open-wallet/internal/resources/categories/repository"
	categoriesUseCases "github.com/felipe1496/open-wallet/internal/resources/categories/usecases"
	recurrencesRepo "github.com/felipe1496/open-wallet/internal/resources/recurrences/repository"
	recurrencesUseCases "github.com/felipe1496/open-wallet/internal/resources/recurrences/usecases"
	transactionsRepo "github.com/felipe1496/open-wallet/internal/resources/transactions/repository"
	transactionsUseCases "github.com/felipe1496/open-wallet/internal/resources/transactions/usecases"
	usersRepo "github.com/felipe1496/open-wallet/internal/resources/users/repository"
	usersUseCases "github.com/felipe1496/open-wallet/internal/resources/users/usecases"
	"github.com/felipe1496/open-wallet/internal/services"
	"github.com/felipe1496/open-wallet/infra/rabbitmq"
)

type Factory struct {
	db  *sql.DB
	cfg *infra.Config

	googleService        services.GoogleService
	jwtService           services.JWTService
	cacheService         services.CacheService
	usersUseCases        usersUseCases.UsersUseCases
	authUseCases         authUseCases.AuthUseCases
	categoriesUseCases   categoriesUseCases.CategoriesUseCases
	transactionsUseCases transactionsUseCases.TransactionsUseCases
	recurrencesUseCases  recurrencesUseCases.RecurrencesUseCases
	brokerService        services.BrokerService
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

func (f *Factory) CacheService() services.CacheService {
	if f.cacheService == nil {
		f.cacheService = services.NewCacheService(f.db)
	}
	return f.cacheService
}

func (f *Factory) UsersUseCases() usersUseCases.UsersUseCases {
	if f.usersUseCases == nil {
		f.usersUseCases = usersUseCases.NewUsersUseCases(usersRepo.NewUsersRepo(), f.db)
	}
	return f.usersUseCases
}

func (f *Factory) AuthUseCases() authUseCases.AuthUseCases {
	if f.authUseCases == nil {
		f.authUseCases = authUseCases.NewAuthUseCases(f.GoogleService(), f.UsersUseCases())
	}
	return f.authUseCases
}

func (f *Factory) CategoriesUseCases() categoriesUseCases.CategoriesUseCases {
	if f.categoriesUseCases == nil {
		f.categoriesUseCases = categoriesUseCases.NewCategoriesUseCases(categoriesRepo.NewCategoriesRepo(), f.db)
	}
	return f.categoriesUseCases
}

func (f *Factory) TransactionsUseCases() transactionsUseCases.TransactionsUseCases {
	if f.transactionsUseCases == nil {
		f.transactionsUseCases = transactionsUseCases.NewTransactionsUseCases(
			transactionsRepo.NewTransactionsRepo(),
			transactionsRepo.NewEntriesRepo(),
			transactionsRepo.NewSummariesRepo(),
			f.CategoriesUseCases(),
			f.db,
		)
	}
	return f.transactionsUseCases
}

func (f *Factory) RecurrencesUseCases() recurrencesUseCases.RecurrencesUseCases {
	if f.recurrencesUseCases == nil {
		f.recurrencesUseCases = recurrencesUseCases.NewRecurrencesUseCases(
			recurrencesRepo.NewRecurrencesRepo(),
			f.CategoriesUseCases(),
			f.TransactionsUseCases(),
			f.db,
		)
	}
	return f.recurrencesUseCases
}

func (f *Factory) BrokerService() (services.BrokerService, error) {
	if f.brokerService == nil {
		client, err := rabbitmq.NewClient(f.cfg.RabbitMQURL)
		if err != nil {
			return nil, err
		}
		f.brokerService = client
	}
	return f.brokerService, nil
}
