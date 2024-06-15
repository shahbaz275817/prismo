package user

import (
	"context"

	"gorm.io/gorm"

	"github.com/shahbaz275817/prismo/internal/models"
	"github.com/shahbaz275817/prismo/internal/repository"
	"github.com/shahbaz275817/prismo/pkg/errors"
)

type Repository interface {
	Get(ctx context.Context, query *models.User) (*models.User, error)
	Save(ctx context.Context, query *models.User) error
	Update(ctx context.Context, query *models.User, update map[string]interface{}) error
}

type userRepository struct {
	dB repository.Accessor
}

func NewUserRepository(accessor repository.Accessor) Repository {
	return &userRepository{
		dB: accessor,
	}
}

func (userRepo *userRepository) Get(ctx context.Context, query *models.User) (*models.User, error) {
	var user models.User

	err := userRepo.dB.Transact(ctx, func(ctx context.Context) error {
		return repository.GetTx(ctx).First(&user, query).Error
	})

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, errors.NewUnknownError(err.Error())
	}

	return &user, nil
}

func (userRepo *userRepository) Save(ctx context.Context, model *models.User) error {
	err := userRepo.dB.Transact(ctx, func(ctx context.Context) error {
		return repository.GetTx(ctx).Create(model).Error
	})
	return err
}

func (userRepo *userRepository) Update(ctx context.Context, model *models.User, update map[string]interface{}) error {

	err := userRepo.dB.Transact(ctx, func(ctx context.Context) error {
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
