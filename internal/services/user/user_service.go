package user

import (
	"context"

	"github.com/shahbaz275817/prismo/internal/models"
	"github.com/shahbaz275817/prismo/pkg/cache/inmemory"

	"github.com/shahbaz275817/prismo/internal/repository/user"
)

type Service interface {
	GetByEmail(ctx context.Context, email string) (user *models.User, err error)
	Update(ctx context.Context, user *models.User, update map[string]interface{}) error
}

type userService struct {
	uCache   inmemory.InMemCache
	userRepo user.Repository
}

func NewUserService(userRepo user.Repository, uCache inmemory.InMemCache) Service {
	return &userService{
		userRepo: userRepo,
		uCache:   uCache,
	}
}

func (service *userService) GetByEmail(ctx context.Context, email string) (user *models.User, err error) {
	query := models.User{Email: email}
	user, err = service.userRepo.Get(ctx, &query)

	if err != nil {
		return nil, err
	}

	return user, err
}

func (service *userService) Update(ctx context.Context, user *models.User, update map[string]interface{}) (err error) {
	err = service.userRepo.Update(ctx, user, update)
	if err != nil {
		return err
	}
	return nil
}
