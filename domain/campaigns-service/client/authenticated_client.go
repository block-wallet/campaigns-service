package client

import (
	"bytes"
	"net/http"

	"github.com/block-wallet/campaigns-service/utils/logger"
)

type authenticatedClient struct {
	accessToken string
	headerKey   string
	debug       bool
}

func NewAuthenticatedClient(accessToken string, headerKey string) *authenticatedClient {
	header := "Authorization"
	if headerKey != "" {
		header = headerKey
	}
	return &authenticatedClient{
		accessToken: accessToken,
		headerKey:   header,
	}
}

func (ac *authenticatedClient) WithDebug(withDebug bool) *authenticatedClient {
	ac.debug = withDebug
	return ac
}

func (ac *authenticatedClient) Do(httpRequest *http.Request) (*http.Response, error) {
	httpRequest.Header.Add(ac.headerKey, ac.accessToken)
	if ac.debug {
		ac.debugRequest(httpRequest)
	}
	client := &http.Client{}
	return client.Do(httpRequest)
}

func (ac *authenticatedClient) debugRequest(httpRequest *http.Request) {
	buf := new(bytes.Buffer)
	ioBody, _ := httpRequest.GetBody()
	buf.ReadFrom(ioBody)
	respBytes := buf.String()
	respString := string(respBytes)
	logger.Sugar.Infof("http-request-headers: %v", httpRequest.Header)
	logger.Sugar.Infof("http-request-url: %v", httpRequest.URL)
	logger.Sugar.Infof("http-request-body: %v", respString)
}
