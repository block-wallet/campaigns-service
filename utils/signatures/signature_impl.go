package signatures

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/storyicon/sigverify"
)

type SimpleMessageVerifier struct {
	message string
}

func NewSimpleMessageVerifier(message string) *SimpleMessageVerifier {
	return &SimpleMessageVerifier{
		message: message,
	}
}

func (v *SimpleMessageVerifier) Verify(signature string, accountAddress string) (bool, error) {
	return sigverify.VerifyEllipticCurveHexSignatureEx(
		common.HexToAddress(accountAddress),
		[]byte(v.message),
		signature,
	)
}
