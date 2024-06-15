package transaction

import (
	"context"

	"gorm.io/gorm"

	"github.com/shahbaz275817/prismo/internal/models"
	"github.com/shahbaz275817/prismo/internal/repository"
	"github.com/shahbaz275817/prismo/pkg/errors"
)

type Repository interface {
	Get(ctx context.Context, query *models.Transaction) (*models.Transaction, error)
	GetAllWithCount(ctx context.Context, query *models.Transaction, request repository.FilterRequest) ([]models.Transaction, int64, error)
	Save(ctx context.Context, query *models.Transaction) error
	Update(ctx context.Context, query *models.Transaction, update *models.Transaction) error
	Transact(ctx context.Context, f func(ctx context.Context) error) error
}

type transactionRepository struct {
	dB repository.Accessor
}

func NewTransactionRepository(accessor repository.Accessor) Repository {
	return &transactionRepository{
		dB: accessor,
	}
}

func (repo *transactionRepository) Get(ctx context.Context, query *models.Transaction) (*models.Transaction, error) {
	var transaction models.Transaction

	err := repo.dB.Transact(ctx, func(ctx context.Context) error {
		return repository.GetTx(ctx).Debug().First(&transaction, query).Error
	})

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, errors.NewUnknownError(err.Error())
	}
	return &transaction, nil
}

func (repo *transactionRepository) GetAllWithCount(ctx context.Context, query *models.Transaction, request repository.FilterRequest) ([]models.Transaction, int64, error) {
	var transactions []models.Transaction
	var count int64

	err := repo.dB.Transact(ctx, func(ctx context.Context) error {
		query := repository.GetTx(ctx).Debug().
			Scopes(
				request.CreatedAtRange(), request.Sort(), request.Pagination()).Where(query).Find(&transactions).
			Offset(-1).
			Count(&count)
		err := query.Error
		return err
	})

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return transactions, 0, nil
	}

	if err != nil {
		return nil, 0, errors.NewUnknownError(err.Error())
	}

	return transactions, count, nil
}

func (repo *transactionRepository) Save(ctx context.Context, model *models.Transaction) error {
	err := repo.dB.Transact(ctx, func(ctx context.Context) error {
		return repository.GetTx(ctx).Create(model).Error
	})
	return err
}

func (repo *transactionRepository) Update(ctx context.Context, model *models.Transaction, update *models.Transaction) error {

	err := repo.dB.Transact(ctx, func(ctx context.Context) error {
		return repository.GetTx(ctx).Model(&model).Updates(update).Error
	})

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil
	}

	if err != nil {
		return nil
	}

	return err
}

func (repo *transactionRepository) Transact(ctx context.Context, f func(ctx context.Context) error) error {
	return repo.dB.Transact(ctx, f)
}
