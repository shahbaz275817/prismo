package account

import (
	"context"

	"github.com/shahbaz275817/prismo/internal/models"
	"github.com/shahbaz275817/prismo/internal/repository/account"
	"github.com/shahbaz275817/prismo/pkg/logger"
)

type Service interface {
	Get(ctx context.Context, query *models.Account) (Account *models.Account, err error)
	Create(ctx context.Context, trx models.Account) error
	Update(ctx context.Context, Account *models.Account, update *models.Account) error
}

type accountService struct {
	repo account.Repository
}

func NewAccountService(repo account.Repository) Service {
	return &accountService{
		repo: repo,
	}
}

func (service *accountService) Get(ctx context.Context, query *models.Account) (Account *models.Account, err error) {
	return service.repo.Get(ctx, query)
}

func (service *accountService) Create(ctx context.Context, acc models.Account) error {

	err := service.repo.Save(ctx, &acc)
	if err != nil {
		logger.WithContext(ctx).Errorf("Error while saving Account Error: %s", err.Error())
		return err
	}

	return nil

}

func (service *accountService) Update(ctx context.Context, Account *models.Account, update *models.Account) (err error) {
	return service.repo.Update(ctx, Account, update)
}
