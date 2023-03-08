package httpapi

import (
	"context"
	"fmt"
	"testing"

	"github.com/block-wallet/golang-service-template/domain/model"
	"github.com/block-wallet/golang-service-template/utils/http/mocks"
	"github.com/smartystreets/assertions"
	"github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/mock"
)

func TestApiClientImpl_GetChains(t *testing.T) {
	cases := []struct {
		name           string
		client         *mocks.Client
		expectedResult *[]model.Chain
		expectedError  error
	}{{
		name: "Should return an error when http client returns an error",
		client: func() *mocks.Client {
			client := &mocks.Client{}
			client.
				On("Get",
					mock.AnythingOfType("*context.emptyCtx"),
					mock.AnythingOfType("string"),
					mock.AnythingOfType("map[string]string")).
				Return(nil, fmt.Errorf("some connection error"))
			return client
		}(),
		expectedResult: nil,
		expectedError:  fmt.Errorf("some connection error"),
	}, {
		name: "Should return chain when http client returns something ok",
		client: func() *mocks.Client {
			client := &mocks.Client{}
			client.
				On("Get",
					mock.AnythingOfType("*context.emptyCtx"),
					mock.AnythingOfType("string"),
					mock.AnythingOfType("map[string]string")).
				Return([]byte("[{\"name\":\"Ethereum Mainnet\",\"chain\":\"ETH\",\"network\":\"mainnet\",\"icon\":\"ethereum\",\"rpc\":[\"https://mainnet.infura.io/v3/${INFURA_API_KEY}\",\"wss://mainnet.infura.io/ws/v3/${INFURA_API_KEY}\",\"https://api.mycryptoapi.com/eth\",\"https://cloudflare-eth.com\"],\"faucets\":[],\"nativeCurrency\":{\"name\":\"Ether\",\"symbol\":\"ETH\",\"decimals\":18},\"infoURL\":\"https://ethereum.org\",\"shortName\":\"eth\",\"chainId\":1,\"networkId\":1,\"slip44\":60,\"ens\":{\"registry\":\"0x00000000000C2E074eC69A0dFb2997BA6C7d2e1e\"},\"explorers\":[{\"name\":\"etherscan\",\"url\":\"https://etherscan.io\",\"standard\":\"EIP3091\"}]}]"), nil)
			return client
		}(),
		expectedResult: &[]model.Chain{{
			Name:    "Ethereum Mainnet",
			Chain:   "ETH",
			Network: "mainnet",
			Icon:    "ethereum",
			Rpc:     []string{"https://mainnet.infura.io/v3/${INFURA_API_KEY}", "wss://mainnet.infura.io/ws/v3/${INFURA_API_KEY}", "https://api.mycryptoapi.com/eth", "https://cloudflare-eth.com"},
			Faucet:  []string(nil),
			NativeCurrency: &model.Currency{
				Name:     "Ether",
				Symbol:   "ETH",
				Decimals: 18,
			},
			InfoURL:   "https://ethereum.org",
			ShortName: "eth",
			ChainId:   1,
			NetworkId: 1,
			Ens:       &model.Ens{Registry: "0x00000000000C2E074eC69A0dFb2997BA6C7d2e1e"},
			Explorers: &[]model.Explorer{{
				Name:     "etherscan",
				Url:      "https://etherscan.io",
				Standard: "EIP3091",
			}},
		}},
		expectedError: nil,
	}}

	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()
			convey.Convey("Given an http client", t, func() {
				clientImpl := NewApiClientImpl(c.client, "http", "some.com")
				convey.Convey("When we get the chains", func() {
					ctx := context.Background()
					// Operation
					chains, err := clientImpl.GetChains(ctx)
					convey.Convey("Then we have the chains", func() {
						// Validation
						convey.So(chains, assertions.ShouldResemble, c.expectedResult)
						convey.So(err, assertions.ShouldResemble, c.expectedError)
						if c.client != nil {
							c.client.AssertNumberOfCalls(t, "Get", 1)
						}
					})
				})
			})
		})
	}
}
