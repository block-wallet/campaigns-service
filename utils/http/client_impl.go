package http

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	httpErrors "github.com/block-wallet/campaigns-service/utils/http/errors"
	"github.com/block-wallet/campaigns-service/utils/logger"
	"github.com/block-wallet/campaigns-service/utils/monitoring"

	"github.com/block-wallet/campaigns-service/utils/monitoring/histogram"
)

// ClientImpl implements the Client interface
type ClientImpl struct {
	client              *http.Client
	latencyMetricSender histogram.LatencyMetricSender
}

func NewClientImpl(timeout time.Duration, latencyMetricSender histogram.LatencyMetricSender) *ClientImpl {
	client := &http.Client{
		Timeout: timeout,
	}
	return &ClientImpl{
		client:              client,
		latencyMetricSender: latencyMetricSender,
	}
}

func (c *ClientImpl) Get(ctx context.Context, urlForResource string, headers map[string]string) ([]byte, error) {
	logger.Sugar.WithCtx(ctx).Debugf("Starting to get response for %s", urlForResource)
	httpReq, err := c.generateHTTPRequest(ctx, "GET", urlForResource, headers)
	if err != nil {
		return nil, err
	}

	httpResp, err := c.do(httpReq)
	if err != nil {
		logger.Sugar.WithCtx(ctx).Errorf("Error executing GET request for url: %s, err: %s", urlForResource, err.Error())
		return nil, err
	}

	body, err := c.processResponse(ctx, httpResp, urlForResource)
	if err != nil {
		return nil, err
	}

	logger.Sugar.WithCtx(ctx).Debugf("Finishing to get response for url: %s", urlForResource)
	return body, nil
}
func (c *ClientImpl) do(httpReq *http.Request) (*http.Response, error) {
	start := time.Now()
	httpResp, err := c.client.Do(httpReq)
	end := time.Now()
	labels := map[string]string{
		monitoring.MethodLabel: httpReq.Method,
		monitoring.StatusLabel: c.getStatusCode(httpResp, err),
	}
	c.latencyMetricSender.Send(start, end, labels)
	return httpResp, err
}

func (c *ClientImpl) processResponse(ctx context.Context, httpResp *http.Response, url string) ([]byte, error) {
	body, readBodyErr := c.readBody(ctx, httpResp, url)
	if http.StatusOK <= httpResp.StatusCode && httpResp.StatusCode < http.StatusMultipleChoices {
		return body, readBodyErr
	}

	var err error
	message := fmt.Sprintf("Error getting response for url: %s, response: %v", url, httpResp)
	switch {
	case http.StatusNotFound == httpResp.StatusCode:
		err = httpErrors.NewNotFound(message)
	case http.StatusBadRequest <= httpResp.StatusCode && httpResp.StatusCode < http.StatusInternalServerError:
		err = httpErrors.NewClient(httpResp.StatusCode, message)
	default:
		err = httpErrors.NewServer(httpResp.StatusCode, message)
	}
	logger.Sugar.WithCtx(ctx).Error(err.Error())
	return nil, err
}

func (c *ClientImpl) readBody(ctx context.Context, httpResp *http.Response, url string) ([]byte, error) {
	body, err := ioutil.ReadAll(httpResp.Body)
	if err != nil {
		logger.Sugar.WithCtx(ctx).Errorf("Error reading body %v for url: %s, err: %s", httpResp.Body, url, err.Error())
		return nil, err
	}

	err = httpResp.Body.Close()
	if err != nil {
		logger.Sugar.WithCtx(ctx).Errorf("Error closing body response %v for url: %s, err: %s", httpResp.Body, url, err.Error())
	}
	return body, nil
}

func (c *ClientImpl) generateHTTPRequest(ctx context.Context, method, url string, headers map[string]string) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		logger.Sugar.WithCtx(ctx).Errorf("Error creating %s request for url: %s, err: %s", method, url, err.Error())
		return nil, err
	}

	req.URL.RawQuery = req.URL.Query().Encode()

	for key, value := range headers {
		req.Header.Add(key, value)
	}

	return req, nil
}

func (c *ClientImpl) getStatusCode(httpResp *http.Response, err error) string {
	if err != nil {
		return "client_error"
	}
	return strconv.Itoa(httpResp.StatusCode)
}
