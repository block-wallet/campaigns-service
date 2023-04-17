package campaignsrepository

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/block-wallet/campaigns-service/domain/model"
	"github.com/stretchr/testify/assert"
)

func Test_Query(t *testing.T) {
	id := "abc123"
	chains := []uint32{1, 137}
	tags := []string{"t1", "t2", "t3"}
	fromDate := time.Now()
	toDate := time.Now().AddDate(0, 1, 0)
	status := []model.CampaignStatus{model.STATUS_ACTIVE, model.STATUS_PENDING}
	cases := []struct {
		name                string
		filters             *model.GetCampaignsFilters
		expectedParamsCount int
		expectedFilters     string
	}{
		{
			name:                "should not have query filters",
			filters:             &model.GetCampaignsFilters{},
			expectedParamsCount: 0,
			expectedFilters:     "",
		},
		{
			name:                "should add only id filter",
			filters:             &model.GetCampaignsFilters{Id: &id},
			expectedParamsCount: 1,
			expectedFilters:     "c.id = $1",
		},
		{
			name:                "should add multiple filters",
			filters:             &model.GetCampaignsFilters{ChainIds: &chains, Tags: &tags, ToDate: &toDate, FromDate: &fromDate, Status: &status},
			expectedParamsCount: 9,
			expectedFilters:     "cscf.chain_id IN ($1,$2) AND ctf.tag IN ($3,$4,$5) AND c.status IN ($6,$7) AND c.start_date::date >= $8::date AND stc.end_date::date <= $9::date",
		},
	}

	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			queryBuilder := NewCampaignsQueryBuilder(c.filters)
			query, p := queryBuilder.Query(context.Background())
			assert.NotEmpty(t, query)
			assert.Contains(t, query, c.expectedFilters)
			assert.Len(t, p, c.expectedParamsCount)
			andIndex := strings.Index(query, "AND")
			if c.expectedFilters == "" {
				assert.Equal(t, andIndex, -1)
			} else {
				assert.NotEqual(t, andIndex, -1)
				filters := query[andIndex:]
				fmt.Println(filters)
				assert.Contains(t, filters, c.expectedFilters)
			}
		})
	}
}
