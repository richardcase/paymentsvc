// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"
import model "github.com/richardcase/paymentsvc/pkg/model"

// Repository is an autogenerated mock type for the Repository type
type Repository struct {
	mock.Mock
}

// Create provides a mock function with given fields: attributes
func (_m *Repository) Create(attributes *model.PaymentAttributes) (*model.Payment, error) {
	ret := _m.Called(attributes)

	var r0 *model.Payment
	if rf, ok := ret.Get(0).(func(*model.PaymentAttributes) *model.Payment); ok {
		r0 = rf(attributes)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Payment)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*model.PaymentAttributes) error); ok {
		r1 = rf(attributes)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Delete provides a mock function with given fields: id
func (_m *Repository) Delete(id string) error {
	ret := _m.Called(id)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetAll provides a mock function with given fields:
func (_m *Repository) GetAll() ([]*model.Payment, error) {
	ret := _m.Called()

	var r0 []*model.Payment
	if rf, ok := ret.Get(0).(func() []*model.Payment); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.Payment)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetByID provides a mock function with given fields: id
func (_m *Repository) GetByID(id string) (*model.Payment, error) {
	ret := _m.Called(id)

	var r0 *model.Payment
	if rf, ok := ret.Get(0).(func(string) *model.Payment); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Payment)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Update provides a mock function with given fields: id, attributes
func (_m *Repository) Update(id string, attributes *model.PaymentAttributes) (*model.Payment, error) {
	ret := _m.Called(id, attributes)

	var r0 *model.Payment
	if rf, ok := ret.Get(0).(func(string, *model.PaymentAttributes) *model.Payment); ok {
		r0 = rf(id, attributes)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Payment)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, *model.PaymentAttributes) error); ok {
		r1 = rf(id, attributes)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
