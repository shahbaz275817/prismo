package account

import (
	"context"

	"gorm.io/gorm"

	"github.com/shahbaz275817/prismo/internal/models"
	"github.com/shahbaz275817/prismo/internal/repository"
	"github.com/shahbaz275817/prismo/pkg/errors"
)

type Repository interface {
	Get(ctx context.Context, query *models.Account) (*models.Account, error)
	Save(ctx context.Context, query *models.Account) error
	Update(ctx context.Context, query *models.Account, update *models.Account) error
}

type accountRepository struct {
	dB repository.Accessor
}

func NewAccountRepository(accessor repository.Accessor) Repository {
	return &accountRepository{
		dB: accessor,
	}
}

func (repo *accountRepository) Get(ctx context.Context, query *models.Account) (*models.Account, error) {
	var Account models.Account

	err := repo.dB.Transact(ctx, func(ctx context.Context) error {
		return repository.GetTx(ctx).Debug().First(&Account, query).Error
	})

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, errors.NewUnknownError(err.Error())
	}
	return &Account, nil
}

func (repo *accountRepository) Save(ctx context.Context, model *models.Account) error {
	err := repo.dB.Transact(ctx, func(ctx context.Context) error {
		return repository.GetTx(ctx).Create(model).Error
	})
	return err
}

func (repo *accountRepository) Update(ctx context.Context, model *models.Account, update *models.Account) error {

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

func (repo *accountRepository) Transact(ctx context.Context, f func(ctx context.Context) error) error {
	return repo.dB.Transact(ctx, f)
}
