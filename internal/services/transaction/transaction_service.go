package transaction

import (
	"context"

	"github.com/shahbaz275817/prismo/internal/models"
	"github.com/shahbaz275817/prismo/internal/repository"
	"github.com/shahbaz275817/prismo/internal/repository/transaction"
	"github.com/shahbaz275817/prismo/pkg/logger"
)

type Service interface {
	Get(ctx context.Context, query *models.Transaction) (transaction *models.Transaction, err error)
	GetAllWithCount(ctx context.Context, query *models.Transaction, request repository.FilterRequest) ([]models.Transaction, int64, error)
	Create(ctx context.Context, trx models.Transaction) (transaction *models.Transaction, err error)
	Update(ctx context.Context, transaction *models.Transaction, update *models.Transaction) error
	Transact(ctx context.Context, f func(ctx context.Context) error) error
}

type transactionService struct {
	repo transaction.Repository
}

func NewTransactionService(repo transaction.Repository) Service {
	return &transactionService{
		repo: repo,
	}
}

func (service *transactionService) Get(ctx context.Context, query *models.Transaction) (transaction *models.Transaction, err error) {
	return service.repo.Get(ctx, query)
}

func (service *transactionService) GetAllWithCount(ctx context.Context, query *models.Transaction, request repository.FilterRequest) ([]models.Transaction, int64, error) {
	return service.repo.GetAllWithCount(ctx, query, request)
}

func (service *transactionService) Create(ctx context.Context, trx models.Transaction) (*models.Transaction, error) {

	err := service.repo.Transact(ctx, func(ctx context.Context) error {
		err := service.repo.Save(ctx, &trx)
		if err != nil {
			logger.WithContext(ctx).Errorf("Error while saving transaction Error: %s", err.Error())
			return err
		}
		return nil
	})

	if err != nil {
		logger.WithContext(ctx).Errorf("Error occurred during create transaction Error: %s", err.Error())
		return nil, err
	}

	return &trx, err
}

func (service *transactionService) Update(ctx context.Context, transaction *models.Transaction, update *models.Transaction) (err error) {
	return service.repo.Update(ctx, transaction, update)
}

func (service *transactionService) Transact(ctx context.Context, f func(ctx context.Context) error) error {
	return service.repo.Transact(ctx, f)
}
