package user

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/shahbaz275817/prismo/internal/config"
	"github.com/shahbaz275817/prismo/internal/models"
	"github.com/shahbaz275817/prismo/internal/repository"
)

type UserRepositoryTestSuit struct {
	suite.Suite
	dB        repository.Accessor
	ctx       context.Context
	userModel *models.User
	userRepo  Repository
}

func (suite *UserRepositoryTestSuit) SetupSuite() {
	config.Load()
	dbAccessor, err := repository.NewAccessor(config.DB())
	if err != nil {
		panic(any("error creating db accessor in consumer repository"))
	}
	suite.dB = dbAccessor
	suite.userRepo = NewUserRepository(suite.dB)
	suite.ctx = context.Background()

	user := &models.User{
		HubID:    123,
		Name:     "Test User",
		Email:    "testuser@gojek.com",
		IsActive: true,
	}
	suite.userModel = user
	err = suite.dB.Transact(suite.ctx, func(ctx context.Context) error {
		return repository.GetTx(ctx).Create(user).Error
	})
	suite.Require().NoError(err)
}

func (suite *UserRepositoryTestSuit) TearDownSuite() {
	err := suite.dB.Transact(suite.ctx, func(ctx context.Context) error {
		return repository.GetTx(ctx).Exec("Truncate table users").Error
	})
	suite.Require().NoError(err)
}

func (suite *UserRepositoryTestSuit) TestGetUserByEmailID() {
	query := &models.User{
		Email: "testuser@gojek.com",
	}
	user, err := suite.userRepo.Get(suite.ctx, query)

	suite.Require().NoError(err)
	suite.Require().EqualValues(suite.userModel.Email, user.Email)
	suite.Require().EqualValues(suite.userModel.Name, user.Name)
	suite.Require().EqualValues(suite.userModel.HubID, user.HubID)
}

func (suite *UserRepositoryTestSuit) TestUpdateUser() {
	query := &models.User{
		Email: "testuser@gojek.com",
	}

	update := map[string]interface{}{
		"HubID":    234,
		"IsActive": false,
	}
	err := suite.userRepo.Update(suite.ctx, query, update)
	suite.Require().NoError(err)

	user, err := suite.userRepo.Get(suite.ctx, query)
	suite.Require().NoError(err)

	suite.Require().EqualValues(suite.userModel.Email, user.Email)
	suite.Require().EqualValues(false, user.IsActive)
	suite.Require().EqualValues(234, user.HubID)
}

func (suite *UserRepositoryTestSuit) TestNotFoundDoesnotReturnErr() {
	query := &models.User{
		Email: "wrong@gojek.com",
	}

	user, err := suite.userRepo.Get(suite.ctx, query)
	suite.Require().NoError(err)
	suite.Assert().Nil(user)
}

func TestUserRepositoryTestSuit(t *testing.T) {
	suite.Run(t, new(UserRepositoryTestSuit))
}
