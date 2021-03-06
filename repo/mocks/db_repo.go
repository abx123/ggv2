// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package mocks

import (
	context "context"
	entities "ggv2/entities"

	mock "github.com/stretchr/testify/mock"
)

// DbRepo is an autogenerated mock type for the DbRepo type
type DbRepo struct {
	mock.Mock
}

// AddToGuestList provides a mock function with given fields: _a0, _a1
func (_m *DbRepo) AddToGuestList(_a0 context.Context, _a1 *entities.Guest) error {
	ret := _m.Called(_a0, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *entities.Guest) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// CreateTable provides a mock function with given fields: _a0, _a1
func (_m *DbRepo) CreateTable(_a0 context.Context, _a1 *entities.Table) (*entities.Table, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *entities.Table
	if rf, ok := ret.Get(0).(func(context.Context, *entities.Table) *entities.Table); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entities.Table)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *entities.Table) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// EmptyTables provides a mock function with given fields: _a0
func (_m *DbRepo) EmptyTables(_a0 context.Context) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetEmptySeatsCount provides a mock function with given fields: _a0
func (_m *DbRepo) GetEmptySeatsCount(_a0 context.Context) (int, error) {
	ret := _m.Called(_a0)

	var r0 int
	if rf, ok := ret.Get(0).(func(context.Context) int); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(int)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetGuestByName provides a mock function with given fields: _a0, _a1
func (_m *DbRepo) GetGuestByName(_a0 context.Context, _a1 *entities.Guest) (*entities.Guest, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *entities.Guest
	if rf, ok := ret.Get(0).(func(context.Context, *entities.Guest) *entities.Guest); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entities.Guest)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *entities.Guest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetTable provides a mock function with given fields: _a0, _a1
func (_m *DbRepo) GetTable(_a0 context.Context, _a1 int64) (*entities.Table, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *entities.Table
	if rf, ok := ret.Get(0).(func(context.Context, int64) *entities.Table); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entities.Table)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int64) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GuestArrived provides a mock function with given fields: _a0, _a1
func (_m *DbRepo) GuestArrived(_a0 context.Context, _a1 *entities.Guest) error {
	ret := _m.Called(_a0, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *entities.Guest) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GuestDepart provides a mock function with given fields: _a0, _a1
func (_m *DbRepo) GuestDepart(_a0 context.Context, _a1 *entities.Guest) error {
	ret := _m.Called(_a0, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *entities.Guest) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ListArrivedGuests provides a mock function with given fields: _a0, _a1, _a2
func (_m *DbRepo) ListArrivedGuests(_a0 context.Context, _a1 int64, _a2 int64) ([]*entities.Guest, error) {
	ret := _m.Called(_a0, _a1, _a2)

	var r0 []*entities.Guest
	if rf, ok := ret.Get(0).(func(context.Context, int64, int64) []*entities.Guest); ok {
		r0 = rf(_a0, _a1, _a2)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*entities.Guest)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int64, int64) error); ok {
		r1 = rf(_a0, _a1, _a2)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListGuests provides a mock function with given fields: _a0, _a1, _a2
func (_m *DbRepo) ListGuests(_a0 context.Context, _a1 int64, _a2 int64) ([]*entities.Guest, error) {
	ret := _m.Called(_a0, _a1, _a2)

	var r0 []*entities.Guest
	if rf, ok := ret.Get(0).(func(context.Context, int64, int64) []*entities.Guest); ok {
		r0 = rf(_a0, _a1, _a2)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*entities.Guest)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int64, int64) error); ok {
		r1 = rf(_a0, _a1, _a2)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListTables provides a mock function with given fields: _a0, _a1, _a2
func (_m *DbRepo) ListTables(_a0 context.Context, _a1 int64, _a2 int64) ([]*entities.Table, error) {
	ret := _m.Called(_a0, _a1, _a2)

	var r0 []*entities.Table
	if rf, ok := ret.Get(0).(func(context.Context, int64, int64) []*entities.Table); ok {
		r0 = rf(_a0, _a1, _a2)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*entities.Table)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int64, int64) error); ok {
		r1 = rf(_a0, _a1, _a2)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
