package redis_hook

import (
	"context"
	"errors"
	"strings"
	"testing"

	m "github.com/block-wallet/golang-service-template/utils/monitoring/histogram"
	"github.com/go-redis/redis/v8"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
)

func TestHook(t *testing.T) {
	assert := assert.New(t)
	n := "namespace"
	p := m.NewPrometheusLatencyMetricSender(n, "redis", "help", []float64{5, 10, 25, 50, 100, 250, 500, 1000, 2500, 5000, 10000}, []string{"method", "status"})

	t.Run("create a new hook", func(t *testing.T) {
		// act
		sut := NewMetricsHook(nil)

		// assert
		assert.NotNil(sut)
	})

	t.Run("do not panic if metrics are already registered", func(t *testing.T) {
		// arrange
		_ = NewMetricsHook(nil)

		// act/assert
		assert.NotPanics(func() {
			_ = NewMetricsHook(nil)
		})
	})

	t.Run("export metrics after a command is processed OK", func(t *testing.T) {
		// arrange
		sut := NewMetricsHook(p)

		cmd := redis.NewStringCmd(context.TODO(), "get")

		// act
		ctx, err1 := sut.BeforeProcess(context.TODO(), cmd)
		err2 := sut.AfterProcess(ctx, cmd)

		// assert
		assert.Nil(err1)
		assert.Nil(err2)

		cmd = redis.NewStringCmd(context.TODO(), "set")
		cmd.SetErr(errors.New("some error"))

		// act
		ctx, err1 = sut.BeforeProcess(context.TODO(), cmd)
		err2 = sut.AfterProcess(ctx, cmd)

		// assert
		assert.Nil(err1)
		assert.Nil(err2)
		metrics, err := prometheus.DefaultGatherer.Gather()
		assert.Nil(err)
		assert.NotNil(metrics)
		var result []map[string]string

		for _, metric := range metrics {
			if strings.HasPrefix(*metric.Name, n) {
				for _, m := range metric.Metric {
					for _, l := range m.GetLabel() {
						result = append(result, map[string]string{*l.Name: *l.Value})
					}
				}
			}
		}

		assert.Len(result, 4)
		assert.Equal("get", result[0]["method"])
		assert.Equal("ok", result[1]["status"])
		assert.Equal("set", result[2]["method"])
		assert.Equal("server_error", result[3]["status"])

	})
}

func TestGetStatus(t *testing.T) {
	assert.Equal(t, GetStatus(nil), "ok")
	assert.Equal(t, GetStatus(redis.Nil), "not_found")
	assert.Equal(t, GetStatus(errors.New("error")), "server_error")

}
