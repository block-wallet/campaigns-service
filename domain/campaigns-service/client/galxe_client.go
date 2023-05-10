package client

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
)

type PopulateParticipantsInput struct {
	Address      common.Address
	CredentialId string
}

type GalxeClient interface {
	PopulateParticipant(ctx context.Context, input PopulateParticipantsInput) (bool, error)
}
