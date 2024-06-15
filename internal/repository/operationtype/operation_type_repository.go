package operationtype

import (
	"context"

	"gorm.io/gorm"

	"github.com/shahbaz275817/prismo/internal/models"
	"github.com/shahbaz275817/prismo/internal/repository"
	"github.com/shahbaz275817/prismo/pkg/errors"
)

type Repository interface {
	Get(ctx context.Context, query *models.OperationsType) (*models.OperationsType, error)
	Save(ctx context.Context, query *models.OperationsType) error
	Update(ctx context.Context, query *models.OperationsType, update *models.OperationsType) error
}

type operationTypeRepository struct {
	dB repository.Accessor
}

func NewOperationTypeRepository(accessor repository.Accessor) Repository {
	return &operationTypeRepository{
		dB: accessor,
	}
}

func (repo *operationTypeRepository) Get(ctx context.Context, query *models.OperationsType) (*models.OperationsType, error) {
	var OperationType models.OperationsType

	err := repo.dB.Transact(ctx, func(ctx context.Context) error {
		return repository.GetTx(ctx).Debug().First(&OperationType, query).Error
	})

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, errors.NewUnknownError(err.Error())
	}
	return &OperationType, nil
}

func (repo *operationTypeRepository) Save(ctx context.Context, model *models.OperationsType) error {
	err := repo.dB.Transact(ctx, func(ctx context.Context) error {
		return repository.GetTx(ctx).Create(model).Error
	})
	return err
}

func (repo *operationTypeRepository) Update(ctx context.Context, model *models.OperationsType, update *models.OperationsType) error {

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

func (repo *operationTypeRepository) Transact(ctx context.Context, f func(ctx context.Context) error) error {
	return repo.dB.Transact(ctx, f)
}
