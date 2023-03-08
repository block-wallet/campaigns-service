package service

import (
	"context"
	"encoding/hex"
	"strconv"
	"strings"

	"github.com/block-wallet/golang-service-template/domain/model"

	"github.com/block-wallet/golang-service-template/domain/eth/contract"
	"github.com/ethereum/go-ethereum/accounts/abi"

	ethrepository "github.com/block-wallet/golang-service-template/domain/eth-service/repository"

	"github.com/ethereum/go-ethereum/crypto"

	"github.com/block-wallet/golang-service-template/utils/logger"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

const (
	Deposit    = "Deposit"
	Withdrawal = "Withdrawal"
)

var CONTRACT, _ = abi.JSON(strings.NewReader(contract.ContractMetaData.ABI))
var SIGNATURES = [...]string{
	"Deposit(bytes32,uint32,uint256)",
	"Withdrawal(address,bytes32,address,uint256)",
}
var TOPICS = map[string]string{
	crypto.Keccak256Hash([]byte(SIGNATURES[0])).Hex(): SIGNATURES[0],
	crypto.Keccak256Hash([]byte(SIGNATURES[1])).Hex(): SIGNATURES[1],
}

type EventsServiceImpl struct {
	Repository        ethrepository.Repository
	ethClient         *ethclient.Client
	contractAddresses []string
	filterQuery       ethereum.FilterQuery
	chLogs            chan types.Log
}

func NewEventsServiceImpl(Repository ethrepository.Repository, url string, contractAddresses []string) (*EventsServiceImpl, error) {
	ethClient, err := ethclient.Dial(url)

	if err != nil {
		logger.Sugar.Errorf("Error dialing eth client - %s", err.Error())
		return nil, err
	}

	addresses := make([]common.Address, len(contractAddresses))

	for i, address := range contractAddresses {
		addresses[i] = common.HexToAddress(address)
	}

	filterQuery := ethereum.FilterQuery{
		Addresses: addresses,
	}

	return &EventsServiceImpl{Repository: Repository, ethClient: ethClient, contractAddresses: contractAddresses, filterQuery: filterQuery, chLogs: make(chan types.Log)}, nil
}

func (s *EventsServiceImpl) Watch(ctx context.Context) error {
	err := s.initialFetching(ctx)
	if err != nil {
		logger.Sugar.WithCtx(ctx).Errorf("Error initializing service - %s", err.Error())
		return err
	}

	subscription, err := s.ethClient.SubscribeFilterLogs(ctx, s.filterQuery, s.chLogs)

	if err != nil {
		logger.Sugar.WithCtx(ctx).Errorf("Error suscribing for logs - %s", err.Error())
		return err
	}

	for {
		select {
		case err := <-subscription.Err():
			logger.Sugar.WithCtx(ctx).Errorf("Error subscription - %s, recovering...", err.Error())
			s.chLogs = make(chan types.Log)
			return s.Watch(ctx)
		case event := <-s.chLogs:
			eventName, ok := TOPICS[event.Topics[0].Hex()]
			if !ok {
				logger.Sugar.Errorf("Error processing event - invalid event")
				continue
			}
			eventName = eventName[:strings.IndexByte(eventName, '(')]

			var indexed abi.Arguments
			for _, arg := range CONTRACT.Events[eventName].Inputs {
				if arg.Indexed {
					indexed = append(indexed, arg)
				}
			}

			switch eventName {
			case Deposit:
				var deposit contract.ContractDeposit
				err := CONTRACT.UnpackIntoInterface(&deposit, eventName, event.Data)
				if err != nil {
					logger.Sugar.WithCtx(ctx).Errorf("Error unpacking event (%s) - %s", eventName, err.Error())
					continue
				}

				err = abi.ParseTopics(&deposit, indexed, event.Topics[1:])
				if err != nil {
					logger.Sugar.WithCtx(ctx).Errorf("Error parsing event (%s) - %s", eventName, err.Error())
					continue
				}

				logger.Sugar.WithCtx(ctx).Debugf("%s: %s - %s", event.Address.String(), eventName, hex.EncodeToString(deposit.Commitment[:]))

				err = s.Repository.SetEvent(ctx, event.Address.String(), strconv.FormatUint(event.BlockNumber, 10), &model.Event{
					BlockNumber:     event.BlockNumber,
					Commitment:      hex.EncodeToString(deposit.Commitment[:]),
					LeafIndex:       deposit.LeafIndex,
					Timestamp:       deposit.Timestamp.String(),
					TransactionHash: event.TxHash.String(),
				})

				if err != nil {
					logger.Sugar.WithCtx(ctx).Errorf("Error saving the event (%s) - %s", eventName, err.Error())
					continue
				}
			case Withdrawal:
				var withdrawal contract.ContractWithdrawal
				err := CONTRACT.UnpackIntoInterface(&withdrawal, eventName, event.Data)
				if err != nil {
					logger.Sugar.WithCtx(ctx).Errorf("Error unpacking event (%s) - %s", eventName, err.Error())
					continue
				}

				err = abi.ParseTopics(&withdrawal, indexed, event.Topics[1:])
				if err != nil {
					logger.Sugar.WithCtx(ctx).Errorf("Error parsing event (%s) - %s", eventName, err.Error())
					continue
				}

				logger.Sugar.WithCtx(ctx).Debugf("%s: %s - %s", event.Address.String(), eventName, withdrawal.Relayer.String())

				// s.Repository.SetWithdrawal(ctx, pair, eventId, event)
			}
		}
	}
}

func (s *EventsServiceImpl) initialFetching(ctx context.Context) error {
	// Fetch and store the historical data from the blockchain
	return nil
}
