// Code generated by mockery v2.38.0. DO NOT EDIT.

package mocks

import (
	context "context"

	models "github.com/shahbaz275817/prismo/internal/models"
	mock "github.com/stretchr/testify/mock"
)

// MockAccountRepository is an autogenerated mock type for the Repository type
type MockAccountRepository struct {
	mock.Mock
}

// Get provides a mock function with given fields: ctx, query
func (_m *MockAccountRepository) Get(ctx context.Context, query *models.Account) (*models.Account, error) {
	ret := _m.Called(ctx, query)

	if len(ret) == 0 {
		panic("no return value specified for Get")
	}

	var r0 *models.Account
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *models.Account) (*models.Account, error)); ok {
		return rf(ctx, query)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *models.Account) *models.Account); ok {
		r0 = rf(ctx, query)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Account)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *models.Account) error); ok {
		r1 = rf(ctx, query)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Save provides a mock function with given fields: ctx, query
func (_m *MockAccountRepository) Save(ctx context.Context, query *models.Account) error {
	ret := _m.Called(ctx, query)

	if len(ret) == 0 {
		panic("no return value specified for Save")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *models.Account) error); ok {
		r0 = rf(ctx, query)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Update provides a mock function with given fields: ctx, query, update
func (_m *MockAccountRepository) Update(ctx context.Context, query *models.Account, update *models.Account) error {
	ret := _m.Called(ctx, query, update)

	if len(ret) == 0 {
		panic("no return value specified for Update")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *models.Account, *models.Account) error); ok {
		r0 = rf(ctx, query, update)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewMockAccountRepository creates a new instance of MockAccountRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockAccountRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockAccountRepository {
	mock := &MockAccountRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
