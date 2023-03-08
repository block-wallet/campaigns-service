package ethservicevalidator

import (
	ethservicev1service "github.com/block-wallet/golang-service-template/protos/ethservicev1/src/eth"
	"github.com/block-wallet/golang-service-template/utils/errors"
	"github.com/block-wallet/golang-service-template/utils/validator"
)

// RequestValidator implements the Validator interface
type RequestValidator struct {
}

func NewRequestValidator() *RequestValidator {
	return &RequestValidator{}
}

func (r *RequestValidator) ValidateGetEventsRequest(req *ethservicev1service.GetEventsMsg) errors.RichError {
	err := validator.FirstNonValid(
		[]validator.SimpleValidation{
			{
				Parameter: req,
				Validator: validator.NonNil,
				ErrorMsg:  "Request provide field: GetEventsMsg",
			},
			{
				Parameter: req.Pair,
				Validator: validator.StringPresent,
				ErrorMsg:  "Request provide field: pair",
			},
		}...)

	if err != nil {
		return err
	}

	return nil
}
