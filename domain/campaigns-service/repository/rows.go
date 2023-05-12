package campaignsrepository

import "time"

type campaignrow struct {
	id                 string
	name               string
	description        string
	status             string
	startDate          string
	endDate            string
	enrollmentMode     string
	createdAt          time.Time
	updatedAt          time.Time
	tags               *string
	supportedChains    string
	campaignType       string
	tokenId            *string
	tokenName          *string
	tokenSymbol        *string
	decimals           *int64
	amounts            *string
	participants       *[]byte
	rewardId           *string
	rewardType         *string
	enrollMessage      *string
	externalCampaignId *string
}

type tokenrow struct {
	id          string
	name        string
	symbol      string
	description string
	decimals    int
}

type participantJSONRow struct {
	AccountAddress  string `json:"account_address"`
	Position        *int   `json:"position"`
	EarlyEnrollment bool   `json:"early_enrollment"`
}
