package model

import (
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
)

const CampaignTimeFormatLayout = time.RFC3339

type CampaignStatus string
type RewardType string

const (
	STATUS_PENDING   CampaignStatus = "PENDING"
	STATUS_ACTIVE    CampaignStatus = "ACTIVE"
	STATUS_FINISHED  CampaignStatus = "FINISHED"
	STATUS_CANCELLED CampaignStatus = "CANCELLED"
	STATUS_UNKNOWN   CampaignStatus = "UNKNOWN"
)

const (
	SINGLE_REWARD  RewardType = "SINGLE_REWARD"
	PODIUM_REWARD  RewardType = "PODIUM_REWARD"
	DYNAMIC_REWARD RewardType = "DYNAMIC_REWARD"
)

type Campaign struct {
	Id              string //uid
	SupportedChains []uint32
	Name            string
	Description     string
	Status          CampaignStatus
	StartDate       time.Time
	EndDate         time.Time
	Rewards         *Reward
	Accounts        []common.Address
	Winners         []common.Address
	Tags            []string
	EnrollMessage   string
}

type MultichainToken struct {
	Id                string
	Name              string
	Symbol            string
	Decimals          uint8
	ContractAddresses map[string]common.Address
}

type Reward struct {
	Type    RewardType
	Token   *MultichainToken
	Amounts []*big.Int
}

type GetCampaignsFilters struct {
	Id       *string
	Status   *[]CampaignStatus
	FromDate *time.Time
	ToDate   *time.Time
	Tags     *[]string
	ChainIds *[]uint32
}

type CampaignRewardTokenInput struct {
	Id     *string
	Create *MultichainToken
}

type CampaignRewardInput struct {
	Amounts []string
	Token   CampaignRewardTokenInput
	Type    RewardType
}

type CreateCampaignInput struct {
	Name            string
	Description     string
	StartDate       string
	EndDate         string
	Tags            []string
	SupportedChains []uint32
	Status          CampaignStatus
	Rewards         CampaignRewardInput
	EnrollMessage   string
}

type EnrollInCampaignInput struct {
	Adddress   common.Address
	CampaignId string
}

type SetCampaignWinners struct {
	Winners []common.Address
}

type UpdateCampaignInput struct {
	Id      string
	Stauts  *CampaignStatus
	Winners *[]common.Address
}
