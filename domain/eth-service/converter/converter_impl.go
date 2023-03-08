package ethserviceconverter

import (
	"github.com/block-wallet/golang-service-template/domain/model"
	ethservicev1service "github.com/block-wallet/golang-service-template/protos/ethservicev1/src/eth"
	"github.com/block-wallet/golang-service-template/utils/grpc/converter"
)

// ConverterImpl has the responsibility to convert from grpc objects to our domain objects and vice versa.
type ConverterImpl struct {
	grpcConverter *converter.GRPCConverter
}

func NewConverterImpl(grpcConverter *converter.GRPCConverter) *ConverterImpl {
	return &ConverterImpl{grpcConverter: grpcConverter}
}

func (c *ConverterImpl) ConvertFromModelEventToProtoEvent(event *model.Event) (*ethservicev1service.Event, error) {
	return &ethservicev1service.Event{
		BlockNumber:     event.BlockNumber,
		Commitment:      event.Commitment,
		LeafIndex:       event.LeafIndex,
		Timestamp:       event.Timestamp,
		TransactionHash: event.TransactionHash,
	}, nil
}
func (c *ConverterImpl) ConvertFromModelChainToProtoChain(chain *model.Chain) (*ethservicev1service.Chain, error) {
	var nativeCurrency *ethservicev1service.NativeCurrency
	if chain.NativeCurrency != nil {
		nativeCurrency = &ethservicev1service.NativeCurrency{
			Name:     chain.NativeCurrency.Name,
			Symbol:   chain.NativeCurrency.Symbol,
			Decimals: chain.NativeCurrency.Decimals,
		}
	}

	var ens *ethservicev1service.Ens
	if chain.Ens != nil {
		ens = &ethservicev1service.Ens{
			Registry: chain.Ens.Registry,
		}
	}

	var explorers []*ethservicev1service.Explorer
	if chain.Explorers != nil {
		for _, explorer := range *chain.Explorers {
			explorers = append(explorers, &ethservicev1service.Explorer{
				Name:     explorer.Name,
				Url:      explorer.Url,
				Standard: explorer.Standard,
			})
		}
	}

	return &ethservicev1service.Chain{
		Name:           chain.Name,
		Chain:          chain.Chain,
		Network:        chain.Network,
		Icon:           chain.Icon,
		Rpc:            chain.Rpc,
		Faucet:         chain.Faucet,
		NativeCurrency: nativeCurrency,
		InfoUrl:        chain.InfoURL,
		ShortName:      chain.ShortName,
		ChainId:        chain.ChainId,
		NetworkId:      chain.NetworkId,
		Ens:            ens,
		Explorers:      explorers,
	}, nil
}

func (c *ConverterImpl) ConvertFromProtoEventToModelEvent(event *ethservicev1service.Event) (*model.Event, error) {
	return &model.Event{
		BlockNumber:     event.GetBlockNumber(),
		Commitment:      event.GetCommitment(),
		LeafIndex:       event.GetLeafIndex(),
		Timestamp:       event.GetTimestamp(),
		TransactionHash: event.GetTransactionHash(),
	}, nil
}
func (c *ConverterImpl) ConvertFromProtoChainToModelChain(chain *ethservicev1service.Chain) (*model.Chain, error) {
	var currency *model.Currency
	if chain.GetNativeCurrency() != nil {
		currency = &model.Currency{
			Name:     chain.GetNativeCurrency().GetName(),
			Symbol:   chain.GetNativeCurrency().GetSymbol(),
			Decimals: chain.GetNativeCurrency().GetDecimals(),
		}
	}

	var ens *model.Ens
	if chain.GetEns() != nil {
		ens = &model.Ens{
			Registry: chain.GetEns().GetRegistry(),
		}
	}

	var explorers []model.Explorer
	if chain.GetExplorers() != nil {
		for _, explorer := range chain.GetExplorers() {
			explorers = append(explorers, model.Explorer{
				Name:     explorer.GetName(),
				Url:      explorer.GetUrl(),
				Standard: explorer.GetStandard(),
			})
		}
	}
	return &model.Chain{
		Name:           chain.GetName(),
		Chain:          chain.GetChain(),
		Network:        chain.GetNetwork(),
		Icon:           chain.GetIcon(),
		Rpc:            chain.GetRpc(),
		Faucet:         chain.GetFaucet(),
		NativeCurrency: currency,
		InfoURL:        chain.GetInfoUrl(),
		ShortName:      chain.GetShortName(),
		ChainId:        chain.GetChainId(),
		NetworkId:      chain.GetNetworkId(),
		Ens:            ens,
		Explorers:      &explorers,
	}, nil
}
