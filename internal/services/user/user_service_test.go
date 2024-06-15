package user

import (
	"context"
	"testing"

	cache_mocks "github.com/shahbaz275817/prismo/pkg/cache/inmemory/mocks"

	"github.com/stretchr/testify/suite"

	"github.com/shahbaz275817/prismo/internal/models"
	"github.com/shahbaz275817/prismo/internal/repository/user/mocks"
	"github.com/shahbaz275817/prismo/pkg/errors"
)

type UserServiceTestSuite struct {
	suite.Suite
	repo  mocks.MockUserRepository
	cache *cache_mocks.MockInMemCache
	ctx   context.Context
}

func (suite *UserServiceTestSuite) SetupTest() {
	suite.ctx = context.Background()
	suite.cache = &cache_mocks.MockInMemCache{}
	suite.repo = mocks.MockUserRepository{}
}

func (suite *UserServiceTestSuite) Test_service_NewUserService() {
	ds := NewUserService(&suite.repo, suite.cache)
	suite.Assert().NotNil(ds)
}

func (suite *UserServiceTestSuite) Test_userService_Update() {
	suite.ctx = context.Background()
	suite.repo = mocks.MockUserRepository{}
	userID := int64(1)
	user := models.User{ID: &userID, HubID: 1}
	type args struct {
		model  *models.User
		update map[string]interface{}
	}
	tests := []struct {
		name     string
		args     args
		mockFunc func()
		wantErr  bool
	}{
		{
			name: "When update success it return nil",
			args: args{
				model:  &user,
				update: map[string]interface{}{"Email": "test1@ad.com"},
			},
			mockFunc: func() {
				suite.repo.On("Update", suite.ctx, &user, map[string]interface{}{"Email": "test1@ad.com"}).Return(nil).Once()
			},
			wantErr: false,
		},
		{
			name: "When update fail it return error",
			args: args{
				model:  &user,
				update: map[string]interface{}{"Email": "test1@ad.com"},
			},
			mockFunc: func() {
				suite.repo.On("Update", suite.ctx, &user, map[string]interface{}{"HubID": int64(2)}).Return(errors.New("something error")).Once()
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		suite.T().Run(tt.name, func(t *testing.T) {
			service := &userService{
				userRepo: &suite.repo,
			}
			tt.mockFunc()
			if err := service.Update(suite.ctx, tt.args.model, tt.args.update); (err != nil) != tt.wantErr {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUserServiceTestSuite(t *testing.T) {
	suite.Run(t, new(UserServiceTestSuite))
}
