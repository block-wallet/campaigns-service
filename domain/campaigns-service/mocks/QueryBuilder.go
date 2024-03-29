// Code generated by mockery v2.23.2. DO NOT EDIT.

package campaignsservicemocks

import (
	context "context"
	sql "database/sql"

	mock "github.com/stretchr/testify/mock"
)

// QueryBuilder is an autogenerated mock type for the QueryBuilder type
type QueryBuilder[P interface{}, F interface{}] struct {
	mock.Mock
}

// Parse provides a mock function with given fields: ctx, row
func (_m *QueryBuilder[P, F]) Parse(ctx context.Context, row *sql.Rows) (*P, error) {
	ret := _m.Called(ctx, row)

	var r0 *P
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *sql.Rows) (*P, error)); ok {
		return rf(ctx, row)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *sql.Rows) *P); ok {
		r0 = rf(ctx, row)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*P)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *sql.Rows) error); ok {
		r1 = rf(ctx, row)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Query provides a mock function with given fields: ctx
func (_m *QueryBuilder[P, F]) Query(ctx context.Context) (string, []string) {
	ret := _m.Called(ctx)

	var r0 string
	var r1 []string
	if rf, ok := ret.Get(0).(func(context.Context) (string, []string)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) string); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(context.Context) []string); ok {
		r1 = rf(ctx)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).([]string)
		}
	}

	return r0, r1
}

type mockConstructorTestingTNewQueryBuilder interface {
	mock.TestingT
	Cleanup(func())
}

// NewQueryBuilder creates a new instance of QueryBuilder. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewQueryBuilder[P interface{}, F interface{}](t mockConstructorTestingTNewQueryBuilder) *QueryBuilder[P, F] {
	mock := &QueryBuilder[P, F]{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
