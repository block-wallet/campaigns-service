package signatures

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Verify(t *testing.T) {
	message := "Campaign service signature test"
	verifier := NewSimpleMessageVerifier(message)
	cases := []struct {
		name           string
		signatureHash  string
		accountAddress string
		isOk           bool
	}{
		{
			name:           "should throw error if signature is not valid",
			signatureHash:  "0xb75654da7683ca5c1b1ab73701c16997991d878904fe8b33a7934a2d81fc381f0df338d2948a4ec3947df0d865d62ee7a366b7e794263e65810461a471f174911c",
			accountAddress: "0x3eaF496818D1a2Cc8aB64C28Df0D2b39c2763211",
			isOk:           false,
		},
		{
			name:           "should validate that the signature is valid",
			accountAddress: "0x3eaF496818D1a2Cc8aB64C28Df0D2b39c2763211",
			signatureHash:  "0x8412d1b6b0b163529b395a49b702ec76a9c1b02a6d46a7031586926575c5b20d4e2896296aa86d976bfb260e85bfe32ac0310f7560e3fdeb02b78a77840779a11b",
			isOk:           true,
		},
	}

	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()
			isOk, _ := verifier.Verify(c.signatureHash, c.accountAddress)
			assert.Equal(t, c.isOk, isOk)
		})
	}
}
