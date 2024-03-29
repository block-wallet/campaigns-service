// Code generated by mockery v2.23.2. DO NOT EDIT.

package campaignsservicemocks

import (
	context "context"

	client "github.com/block-wallet/campaigns-service/domain/campaigns-service/client"

	mock "github.com/stretchr/testify/mock"
)

// GalxeClient is an autogenerated mock type for the GalxeClient type
type GalxeClient struct {
	mock.Mock
}

// PopulateParticipant provides a mock function with given fields: ctx, input
func (_m *GalxeClient) PopulateParticipant(ctx context.Context, input client.PopulateParticipantsInput) (bool, error) {
	ret := _m.Called(ctx, input)

	var r0 bool
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, client.PopulateParticipantsInput) (bool, error)); ok {
		return rf(ctx, input)
	}
	if rf, ok := ret.Get(0).(func(context.Context, client.PopulateParticipantsInput) bool); ok {
		r0 = rf(ctx, input)
	} else {
		r0 = ret.Get(0).(bool)
	}

	if rf, ok := ret.Get(1).(func(context.Context, client.PopulateParticipantsInput) error); ok {
		r1 = rf(ctx, input)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewGalxeClient interface {
	mock.TestingT
	Cleanup(func())
}

// NewGalxeClient creates a new instance of GalxeClient. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewGalxeClient(t mockConstructorTestingTNewGalxeClient) *GalxeClient {
	mock := &GalxeClient{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
