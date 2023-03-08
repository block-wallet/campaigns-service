package ethserviceconverter

import (
	"github.com/block-wallet/golang-service-template/domain/model"
	ethservicev1service "github.com/block-wallet/golang-service-template/protos/ethservicev1/src/eth"
)

type Converter interface {
	ConvertFromModelEventToProtoEvent(*model.Event) (*ethservicev1service.Event, error)
	ConvertFromModelChainToProtoChain(*model.Chain) (*ethservicev1service.Chain, error)

	ConvertFromProtoEventToModelEvent(*ethservicev1service.Event) (*model.Event, error)
	ConvertFromProtoChainToModelChain(*ethservicev1service.Chain) (*model.Chain, error)
}
