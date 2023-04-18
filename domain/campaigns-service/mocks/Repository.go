// Code generated by mockery v2.23.2. DO NOT EDIT.

package campaignsservicemocks

import (
	context "context"

	model "github.com/block-wallet/campaigns-service/domain/model"
	mock "github.com/stretchr/testify/mock"
)

// Repository is an autogenerated mock type for the Repository type
type Repository struct {
	mock.Mock
}

// EnrollInCampaign provides a mock function with given fields: ctx, input
func (_m *Repository) EnrollInCampaign(ctx context.Context, input *model.EnrollInCampaignInput) (*bool, error) {
	ret := _m.Called(ctx, input)

	var r0 *bool
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *model.EnrollInCampaignInput) (*bool, error)); ok {
		return rf(ctx, input)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *model.EnrollInCampaignInput) *bool); ok {
		r0 = rf(ctx, input)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*bool)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *model.EnrollInCampaignInput) error); ok {
		r1 = rf(ctx, input)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAllTokens provides a mock function with given fields: ctx
func (_m *Repository) GetAllTokens(ctx context.Context) ([]*model.MultichainToken, error) {
	ret := _m.Called(ctx)

	var r0 []*model.MultichainToken
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) ([]*model.MultichainToken, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) []*model.MultichainToken); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.MultichainToken)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetCampaignById provides a mock function with given fields: ctx, id
func (_m *Repository) GetCampaignById(ctx context.Context, id string) (*model.Campaign, error) {
	ret := _m.Called(ctx, id)

	var r0 *model.Campaign
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (*model.Campaign, error)); ok {
		return rf(ctx, id)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) *model.Campaign); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Campaign)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetCampaigns provides a mock function with given fields: ctx, filters
func (_m *Repository) GetCampaigns(ctx context.Context, filters *model.GetCampaignsFilters) ([]*model.Campaign, error) {
	ret := _m.Called(ctx, filters)

	var r0 []*model.Campaign
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *model.GetCampaignsFilters) ([]*model.Campaign, error)); ok {
		return rf(ctx, filters)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *model.GetCampaignsFilters) []*model.Campaign); ok {
		r0 = rf(ctx, filters)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.Campaign)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *model.GetCampaignsFilters) error); ok {
		r1 = rf(ctx, filters)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetTokenById provides a mock function with given fields: ctx, id
func (_m *Repository) GetTokenById(ctx context.Context, id string) (*model.MultichainToken, error) {
	ret := _m.Called(ctx, id)

	var r0 *model.MultichainToken
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (*model.MultichainToken, error)); ok {
		return rf(ctx, id)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) *model.MultichainToken); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.MultichainToken)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewCampaign provides a mock function with given fields: ctx, input
func (_m *Repository) NewCampaign(ctx context.Context, input *model.CreateCampaignInput) (*string, error) {
	ret := _m.Called(ctx, input)

	var r0 *string
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *model.CreateCampaignInput) (*string, error)); ok {
		return rf(ctx, input)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *model.CreateCampaignInput) *string); ok {
		r0 = rf(ctx, input)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*string)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *model.CreateCampaignInput) error); ok {
		r1 = rf(ctx, input)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewToken provides a mock function with given fields: ctx, token
func (_m *Repository) NewToken(ctx context.Context, token *model.MultichainToken) (*string, error) {
	ret := _m.Called(ctx, token)

	var r0 *string
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *model.MultichainToken) (*string, error)); ok {
		return rf(ctx, token)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *model.MultichainToken) *string); ok {
		r0 = rf(ctx, token)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*string)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *model.MultichainToken) error); ok {
		r1 = rf(ctx, token)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ParticipantExists provides a mock function with given fields: ctx, campaignId, accountAddress
func (_m *Repository) ParticipantExists(ctx context.Context, campaignId string, accountAddress string) (*bool, error) {
	ret := _m.Called(ctx, campaignId, accountAddress)

	var r0 *bool
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) (*bool, error)); ok {
		return rf(ctx, campaignId, accountAddress)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, string) *bool); ok {
		r0 = rf(ctx, campaignId, accountAddress)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*bool)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(ctx, campaignId, accountAddress)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// TokenExists provides a mock function with given fields: ctx, id
func (_m *Repository) TokenExists(ctx context.Context, id string) (*bool, error) {
	ret := _m.Called(ctx, id)

	var r0 *bool
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (*bool, error)); ok {
		return rf(ctx, id)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) *bool); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*bool)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateCampaign provides a mock function with given fields: ctx, updates
func (_m *Repository) UpdateCampaign(ctx context.Context, updates *model.UpdateCampaignInput) (*bool, error) {
	ret := _m.Called(ctx, updates)

	var r0 *bool
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *model.UpdateCampaignInput) (*bool, error)); ok {
		return rf(ctx, updates)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *model.UpdateCampaignInput) *bool); ok {
		r0 = rf(ctx, updates)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*bool)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *model.UpdateCampaignInput) error); ok {
		r1 = rf(ctx, updates)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewRepository interface {
	mock.TestingT
	Cleanup(func())
}

// NewRepository creates a new instance of Repository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewRepository(t mockConstructorTestingTNewRepository) *Repository {
	mock := &Repository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
