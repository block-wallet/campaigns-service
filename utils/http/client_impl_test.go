package http

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/block-wallet/campaigns-service/utils/http/mocks"
	mocks3 "github.com/block-wallet/campaigns-service/utils/monitoring/histogram/mocks"
	"github.com/smartystreets/assertions"
	"github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/mock"
)

func TestGenerateHTTPRequestWithErr(t *testing.T) {
	convey.Convey("Given a nil context", t, func() {
		var ctx context.Context
		convey.Convey("When I ask for a new server error", func() {
			clientImpl := NewClientImpl(0, nil)
			httpReq, err := clientImpl.generateHTTPRequest(ctx, "", "", nil)
			convey.Convey("Then httpReq should be nil and err should not be nil", func() {
				convey.So(clientImpl, assertions.ShouldNotBeNil)
				convey.So(httpReq, assertions.ShouldBeNil)
				convey.So(err, assertions.ShouldNotBeNil)
				convey.So(err.Error(), assertions.ShouldContainSubstring, "nil Context")
			})
		})
	})
}

func TestGenerateHTTPRequestWithoutErr(t *testing.T) {
	cases := []struct {
		name    string
		headers map[string]string
	}{
		{
			"Should return http req without headers",
			nil,
		},
		{
			"Should return http req without headers",
			map[string]string{},
		},
		{
			"Should return http req with one header",
			map[string]string{
				"header1": "value1",
			},
		},
		{
			"Should return http req with two headers",
			map[string]string{
				"header2": "value2",
			},
		},
	}

	clientImpl := NewClientImpl(0, nil)
	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()
			method := "GET"
			url := "http://localhost:8080?b_key=b_value&a_key=a_value1 a_value2"
			convey.Convey("When we generate the ", t, func() {
				httpReq, err := clientImpl.generateHTTPRequest(ctx, method, url, c.headers)
				convey.Convey("Then we have to get the httpReq with the attributes", func() {
					convey.So(httpReq, assertions.ShouldNotBeNil)
					convey.So(httpReq.Context(), assertions.ShouldResemble, ctx)
					convey.So(httpReq.Method, assertions.ShouldEqual, method)
					convey.So(httpReq.URL.String(), assertions.ShouldEqual, "http://localhost:8080?a_key=a_value1+a_value2&b_key=b_value")
					for key, value := range c.headers {
						convey.So(httpReq.Header.Get(key), assertions.ShouldEqual, value)
					}
					convey.So(err, assertions.ShouldBeNil)
				})
			})
		})
	}
}

func TestReadBodyWithErr(t *testing.T) {
	convey.Convey("Given a httpResp with mocked body", t, func() {
		clientImpl := NewClientImpl(0, nil)
		readCloserMock := &mocks.ReadCloser{}
		readCloserMock.On("Read", mock.Anything).
			Return(0, fmt.Errorf("error reading"))
		httpResp := &http.Response{
			Body: readCloserMock,
		}
		convey.Convey("When we read the body", func() {
			body, err := clientImpl.readBody(context.Background(), httpResp, "")
			convey.Convey("Then we are going to have an error", func() {
				convey.So(body, assertions.ShouldBeNil)
				convey.So(err, assertions.ShouldNotBeNil)
				convey.So(err.Error(), assertions.ShouldContainSubstring, "error reading")
				readCloserMock.AssertNumberOfCalls(t, "Read", 1)
			})
		})
	})
}

func TestReadBodyWithoutErr(t *testing.T) {
	convey.Convey("Given a httpResp with real body", t, func() {
		clientImpl := NewClientImpl(0, nil)
		httpResp := &http.Response{
			Body: ioutil.NopCloser(bytes.NewReader([]byte("hola"))),
		}
		convey.Convey("When we read the body", func() {
			body, err := clientImpl.readBody(context.Background(), httpResp, "")
			convey.Convey("Then we aren't going to have an error", func() {
				convey.So(body, assertions.ShouldNotBeNil)
				convey.So(string(body), assertions.ShouldEqual, "hola")
				convey.So(err, assertions.ShouldBeNil)
			})
		})
	})
}

func TestGet(t *testing.T) {
	serverWithoutTimeout := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Header.Get("status") {
		case "ok":
			w.WriteHeader(http.StatusOK)
		case "clienterror":
			w.WriteHeader(http.StatusBadRequest)
		case "servererror":
			w.WriteHeader(http.StatusInternalServerError)
		case "notfound":
			w.WriteHeader(http.StatusNotFound)
		}
		_, _ = fmt.Fprint(w, "request without timeout")
	}))
	defer serverWithoutTimeout.Close()

	cases := []struct {
		name              string
		ctx               context.Context
		url               string
		headers           map[string]string
		expectedBody      []byte
		expectedErr       string
		metricLabelValues []string
	}{
		{
			"Should returns an error when the context is nil",
			nil,
			serverWithoutTimeout.URL,
			nil,
			nil,
			"nil Context",
			nil,
		},
		{
			"Should returns an error when the client returns 400",
			context.Background(),
			serverWithoutTimeout.URL,
			map[string]string{"status": "clienterror"},
			nil,
			"HTTP client error",
			[]string{serverWithoutTimeout.URL, "400"},
		},
		{
			"Should returns an error when the client returns 500",
			context.Background(),
			serverWithoutTimeout.URL,
			map[string]string{"status": "servererror"},
			nil,
			"HTTP server error",
			[]string{serverWithoutTimeout.URL, "500"},
		},
		{
			"Should returns a nil error when the client returns 200",
			context.Background(),
			serverWithoutTimeout.URL,
			map[string]string{"status": "ok"},
			[]byte("request without timeout"),
			"",
			[]string{serverWithoutTimeout.URL, "200"},
		},
		{
			"Should returns an error when the client returns 404",
			context.Background(),
			serverWithoutTimeout.URL,
			map[string]string{"status": "notfound"},
			nil,
			"HTTP not found error",
			[]string{serverWithoutTimeout.URL, "404"},
		},
	}

	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			convey.Convey("When we get the response", t, func() {
				latencyMetricSender := mocks3.LatencyMetricSender{}
				if c.metricLabelValues != nil {
					latencyMetricSender.On("Send", mock.AnythingOfType("time.Time"), mock.AnythingOfType("time.Time"),
						map[string]string{"method": "GET", "status": c.metricLabelValues[1]})
				}
				clientImplWithoutTimeout := NewClientImpl(0, &latencyMetricSender)
				body, err := clientImplWithoutTimeout.Get(c.ctx, c.url, c.headers)
				convey.Convey("Then we have to get the body and the error", func() {
					convey.So(string(body), assertions.ShouldEqual, string(c.expectedBody))
					if c.expectedErr == "" {
						convey.So(err, assertions.ShouldBeNil)
					} else {
						convey.So(err, assertions.ShouldNotBeNil)
						convey.So(err.Error(), assertions.ShouldContainSubstring, c.expectedErr)
					}
					if c.metricLabelValues != nil {
						latencyMetricSender.AssertCalled(t, "Send", mock.AnythingOfType("time.Time"), mock.AnythingOfType("time.Time"),
							mock.AnythingOfType("map[string]string"))
					}
				})
			})
		})
	}
}

func TestGet_WithTimeout(t *testing.T) {
	serverWithTimeout := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(1 * time.Millisecond)
		_, _ = fmt.Fprintln(w, "request with timeout")
	}))
	defer serverWithTimeout.Close()
	latencyMetricSender := mocks3.LatencyMetricSender{}
	latencyMetricSender.On("Send", mock.AnythingOfType("time.Time"), mock.AnythingOfType("time.Time"),
		map[string]string{"method": "GET", "status": "client_error"})
	clientImplWithTimeout := NewClientImpl(1*time.Nanosecond, &latencyMetricSender)

	convey.Convey("When we get the response", t, func() {
		body, err := clientImplWithTimeout.Get(context.Background(), serverWithTimeout.URL, nil)
		convey.Convey("Then we have to get the body and the error", func() {
			convey.So(body, assertions.ShouldBeNil)
			convey.So(err, assertions.ShouldNotBeNil)
			convey.So(err.Error(), assertions.ShouldContainSubstring, "Client.Timeout exceeded")
			latencyMetricSender.AssertCalled(t, "Send", mock.AnythingOfType("time.Time"), mock.AnythingOfType("time.Time"),
				mock.AnythingOfType("map[string]string"))
		})
	})

}
