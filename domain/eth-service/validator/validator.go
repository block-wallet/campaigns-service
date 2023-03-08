package ethservicevalidator

import (
	ethservicev1service "github.com/block-wallet/golang-service-template/protos/ethservicev1/src/eth"
	"github.com/block-wallet/golang-service-template/utils/errors"
)

type Validator interface {
	ValidateGetEventsRequest(req *ethservicev1service.GetEventsMsg) errors.RichError
}
