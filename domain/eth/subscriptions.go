package eth

import (
	"context"
	"fmt"
	"strings"

	ethrepository "github.com/block-wallet/golang-service-template/domain/eth-service/repository"

	etheventsseervice "github.com/block-wallet/golang-service-template/domain/eth/service"

	ethservice "github.com/block-wallet/golang-service-template/domain/eth-service/service"
	"github.com/block-wallet/golang-service-template/utils/logger"
)

type Subscriptions struct {
	ETHService  ethservice.Service
	Repository  ethrepository.Repository
	ETHEndpoint string
}

func NewSubscriptions(ethService ethservice.Service, repository ethrepository.Repository, ethEndpoint string) *Subscriptions {
	return &Subscriptions{ETHService: ethService, Repository: repository, ETHEndpoint: ethEndpoint}
}

func (s *Subscriptions) BackgroundStart(ctx context.Context) {
	go s.start(ctx)
}

func (s *Subscriptions) start(ctx context.Context) {
	instances := []string{
		"0x6Bf694a291DF3FeC1f7e69701E3ab6c592435Ae7",
		"0x3aac1cC67c2ec5Db4eA850957b967Ba153aD6279",
		"0x723B78e67497E85279CB204544566F4dC5d2acA0",
		"0x0E3A09dDA6B20aFbB34aC7cD4A6881493f3E7bf7",
		"0x76D85B4C0Fc497EeCc38902397aC608000A06607",
		"0xCC84179FFD19A1627E79F8648d09e095252Bc418",
		"0xD5d6f8D9e784d0e26222ad3834500801a68D027D",
		"0x407CcEeaA7c95d2FE2250Bf9F2c105aA7AAFB512",
		"0x833481186f16Cece3f1Eeea1a694c42034c3a0dB",
		"0xd8D7DE3349ccaA0Fde6298fe6D7b7d0d34586193",
		"0x8281Aa6795aDE17C8973e1aedcA380258Bc124F9",
		"0x57b2B8c82F065de8Ef5573f9730fC1449B403C9f",
		"0x05E0b5B40B7b66098C2161A5EE11C5740A3A7C45",
		"0x23173fE8b96A4Ad8d2E17fB83EA5dcccdCa1Ae52",
		"0x538Ab61E8A9fc1b2f93b3dd9011d662d89bE6FE6",
		"0x94Be88213a387E992Dd87DE56950a9aef34b9448",
		"0x242654336ca2205714071898f67E254EB49ACdCe",
		"0x776198CCF446DFa168347089d7338879273172cF",
		"0xeDC5d01286f99A066559F60a585406f3878a033e",
	}

	if len(instances) < 1 {
		logger.Sugar.WithCtx(ctx).Errorf("Error getting instances")
		panic(fmt.Errorf("error getting instances"))
	}

	service, _err := etheventsseervice.NewEventsServiceImpl(s.Repository, s.ETHEndpoint, instances)

	if _err != nil {
		logger.Sugar.WithCtx(ctx).Errorf("Error creating events service - %s", _err.Error())
		panic(_err)
	}

	go func() {
		for {
			err := service.Watch(ctx)
			if err != nil {
				if strings.Contains(err.Error(), "connection reset by peer") {
					logger.Sugar.WithCtx(ctx).Debugf("Recovering from error - %s", err.Error())
					continue
				}
				if strings.Contains(err.Error(), "i/o timeout") {
					logger.Sugar.WithCtx(ctx).Debugf("Recovering from error - %s", err.Error())
					continue
				}
				if strings.Contains(err.Error(), "connection timed out") {
					logger.Sugar.WithCtx(ctx).Debugf("Recovering from error - %s", err.Error())
					continue
				}
				if strings.Contains(err.Error(), "bad handshake (HTTP status 503 Service Temporarily Unavailable)") {
					logger.Sugar.WithCtx(ctx).Debugf("Recovering from error - %s", err.Error())
					continue
				}

				logger.Sugar.WithCtx(ctx).Errorf("Error watching the eth - %s", err.Error())
				panic(err)
			}
		}
	}()

	logger.Sugar.WithCtx(ctx).Infof("Tornado contracts fetched and being watched")
}
