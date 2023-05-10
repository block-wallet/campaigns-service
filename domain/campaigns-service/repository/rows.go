package campaignsrepository

type campaignrow struct {
	id                 string
	name               string
	description        string
	status             string
	startDate          string
	endDate            string
	enrollmentMode     string
	tags               *string
	supportedChains    string
	campaignType       string
	tokenId            *string
	tokenName          *string
	tokenSymbol        *string
	decimals           *int64
	amounts            *string
	participants       *string
	winners            *string
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
