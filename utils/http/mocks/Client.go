// Code generated by mockery (devel). DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// Client is an autogenerated mock type for the Client type
type Client struct {
	mock.Mock
}

// Get provides a mock function with given fields: ctx, url, headers
func (_m *Client) Get(ctx context.Context, url string, headers map[string]string) ([]byte, error) {
	ret := _m.Called(ctx, url, headers)

	var r0 []byte
	if rf, ok := ret.Get(0).(func(context.Context, string, map[string]string) []byte); ok {
		r0 = rf(ctx, url, headers)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, map[string]string) error); ok {
		r1 = rf(ctx, url, headers)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
