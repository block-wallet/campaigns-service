package model

import (
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
)

const CampaignTimeFormatLayout = time.RFC3339

type CampaignStatus string
type RewardType string
type EnrollmentMode string
type CampaignType string

const (
	STATUS_PENDING   CampaignStatus = "PENDING"
	STATUS_WAITLIST  CampaignStatus = "WAITLIST"
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

const (
	INSTANCE_UNLIMITED_ENROLL EnrollmentMode = "INSTANCE_UNLIMITED_ENROLL"
	INSTANCE_SINGLE_ENROLL    EnrollmentMode = "INSTANCE_SINGLE_ENROLL"
)

const (
	CAMPAIGN_TYPE_PARTNER_OFFERS CampaignType = "PARTNER_OFFERS"
	CAMPAIGN_TYPE_GALXE          CampaignType = "GALXE"
	CAMPAIGN_TYPE_STAKING        CampaignType = "STAKING"
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
	Tags            []string
	EnrollMessage   string
	EnrollmentMode  EnrollmentMode
	Type            CampaignType
	Metadata        CampaignMetadata
	Participants    []*CampaignParticipant
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type CampaignParticipant struct {
	AccountAddress  common.Address
	EarlyEnrollment bool
	Position        *int
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

type GalxeCampaignMetadata struct {
	CredentialId string
}

type PartnerOffersMetadata struct {
}

type CampaignMetadata struct {
	GalxeMetadata         *GalxeCampaignMetadata
	PartnerOffersMetadata *PartnerOffersMetadata
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
	EnrollmentMode  EnrollmentMode
	Type            CampaignType
	Metadata        CampaignMetadata
}

type EnrollInCampaignInput struct {
	Adddress        common.Address
	CampaignId      string
	EarlyEnrollment bool
}

type UnenrollFromCampaignInput struct {
	Adddress   common.Address
	CampaignId string
}

type UpdateCampaignInput struct {
	Id               string
	Status           *CampaignStatus
	EligibleAccounts *[]common.Address
}
