package operationtype

import (
	"context"

	"github.com/shahbaz275817/prismo/internal/models"
	"github.com/shahbaz275817/prismo/internal/repository/operationtype"
	"github.com/shahbaz275817/prismo/pkg/logger"
)

type Service interface {
	Get(ctx context.Context, query *models.OperationsType) (OperationType *models.OperationsType, err error)
	Create(ctx context.Context, trx models.OperationsType) error
	Update(ctx context.Context, OperationType *models.OperationsType, update *models.OperationsType) error
}

type operationTypeService struct {
	repo operationtype.Repository
}

func NewOperationTypeService(repo operationtype.Repository) Service {
	return &operationTypeService{
		repo: repo,
	}
}

func (service *operationTypeService) Get(ctx context.Context, query *models.OperationsType) (OperationType *models.OperationsType, err error) {
	return service.repo.Get(ctx, query)
}

func (service *operationTypeService) Create(ctx context.Context, trx models.OperationsType) error {

	err := service.repo.Save(ctx, &trx)
	if err != nil {
		logger.WithContext(ctx).Errorf("Error while saving OperationsType Error: %s", err.Error())
		return err
	}

	return nil

}

func (service *operationTypeService) Update(ctx context.Context, OperationType *models.OperationsType, update *models.OperationsType) (err error) {
	return service.repo.Update(ctx, OperationType, update)
}
