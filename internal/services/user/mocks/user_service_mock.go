// Code generated by mockery v2.31.1. DO NOT EDIT.

package mocks

import (
	context "context"

	models "github.com/shahbaz275817/prismo/internal/models"
	mock "github.com/stretchr/testify/mock"
)

// MockUserService is an autogenerated mock type for the Service type
type MockUserService struct {
	mock.Mock
}

// GetByEmail provides a mock function with given fields: ctx, email
func (_m *MockUserService) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	ret := _m.Called(ctx, email)

	var r0 *models.User
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (*models.User, error)); ok {
		return rf(ctx, email)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) *models.User); ok {
		r0 = rf(ctx, email)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.User)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, email)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetUserWithHubsByID provides a mock function with given fields: ctx, ID
func (_m *MockUserService) GetUserWithHubsByID(ctx context.Context, ID *int64) (*models.User, error) {
	ret := _m.Called(ctx, ID)

	var r0 *models.User
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *int64) (*models.User, error)); ok {
		return rf(ctx, ID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *int64) *models.User); ok {
		r0 = rf(ctx, ID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.User)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *int64) error); ok {
		r1 = rf(ctx, ID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Update provides a mock function with given fields: ctx, _a1, update
func (_m *MockUserService) Update(ctx context.Context, _a1 *models.User, update map[string]interface{}) error {
	ret := _m.Called(ctx, _a1, update)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *models.User, map[string]interface{}) error); ok {
		r0 = rf(ctx, _a1, update)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdateHubPreference provides a mock function with given fields: ctx, _a1, hubID
func (_m *MockUserService) UpdateHubPreference(ctx context.Context, _a1 *models.User, hubID int64) error {
	ret := _m.Called(ctx, _a1, hubID)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *models.User, int64) error); ok {
		r0 = rf(ctx, _a1, hubID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewMockUserService creates a new instance of MockUserService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockUserService(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockUserService {
	mock := &MockUserService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
