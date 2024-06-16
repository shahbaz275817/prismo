// Code generated by mockery v2.38.0. DO NOT EDIT.

package mocks

import (
	context "context"

	models "github.com/shahbaz275817/prismo/internal/models"
	mock "github.com/stretchr/testify/mock"

	repository "github.com/shahbaz275817/prismo/internal/repository"
)

// MockTransactionRepository is an autogenerated mock type for the Repository type
type MockTransactionRepository struct {
	mock.Mock
}

// Get provides a mock function with given fields: ctx, query
func (_m *MockTransactionRepository) Get(ctx context.Context, query *models.Transaction) (*models.Transaction, error) {
	ret := _m.Called(ctx, query)

	if len(ret) == 0 {
		panic("no return value specified for Get")
	}

	var r0 *models.Transaction
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *models.Transaction) (*models.Transaction, error)); ok {
		return rf(ctx, query)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *models.Transaction) *models.Transaction); ok {
		r0 = rf(ctx, query)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Transaction)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *models.Transaction) error); ok {
		r1 = rf(ctx, query)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAllWithCount provides a mock function with given fields: ctx, query, request
func (_m *MockTransactionRepository) GetAllWithCount(ctx context.Context, query *models.Transaction, request repository.FilterRequest) ([]models.Transaction, int64, error) {
	ret := _m.Called(ctx, query, request)

	if len(ret) == 0 {
		panic("no return value specified for GetAllWithCount")
	}

	var r0 []models.Transaction
	var r1 int64
	var r2 error
	if rf, ok := ret.Get(0).(func(context.Context, *models.Transaction, repository.FilterRequest) ([]models.Transaction, int64, error)); ok {
		return rf(ctx, query, request)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *models.Transaction, repository.FilterRequest) []models.Transaction); ok {
		r0 = rf(ctx, query, request)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]models.Transaction)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *models.Transaction, repository.FilterRequest) int64); ok {
		r1 = rf(ctx, query, request)
	} else {
		r1 = ret.Get(1).(int64)
	}

	if rf, ok := ret.Get(2).(func(context.Context, *models.Transaction, repository.FilterRequest) error); ok {
		r2 = rf(ctx, query, request)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// Save provides a mock function with given fields: ctx, query
func (_m *MockTransactionRepository) Save(ctx context.Context, query *models.Transaction) error {
	ret := _m.Called(ctx, query)

	if len(ret) == 0 {
		panic("no return value specified for Save")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *models.Transaction) error); ok {
		r0 = rf(ctx, query)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Transact provides a mock function with given fields: ctx, f
func (_m *MockTransactionRepository) Transact(ctx context.Context, f func(context.Context) error) error {
	ret := _m.Called(ctx, f)

	if len(ret) == 0 {
		panic("no return value specified for Transact")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, func(context.Context) error) error); ok {
		r0 = rf(ctx, f)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Update provides a mock function with given fields: ctx, query, update
func (_m *MockTransactionRepository) Update(ctx context.Context, query *models.Transaction, update *models.Transaction) error {
	ret := _m.Called(ctx, query, update)

	if len(ret) == 0 {
		panic("no return value specified for Update")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *models.Transaction, *models.Transaction) error); ok {
		r0 = rf(ctx, query, update)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewMockTransactionRepository creates a new instance of MockTransactionRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockTransactionRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockTransactionRepository {
	mock := &MockTransactionRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
