package redis_hook

import (
	"context"
	"time"

	"github.com/block-wallet/golang-service-template/utils/monitoring"
	histogrammonitoring "github.com/block-wallet/golang-service-template/utils/monitoring/histogram"

	"github.com/go-redis/redis/v8"
)

type MetricsHook struct {
	latencyMetricSender histogrammonitoring.LatencyMetricSender
}
type startKey struct{}

func NewMetricsHook(latencyMetricSender histogrammonitoring.LatencyMetricSender) *MetricsHook {
	return &MetricsHook{
		latencyMetricSender: latencyMetricSender,
	}
}
func (m *MetricsHook) BeforeProcess(ctx context.Context, cmd redis.Cmder) (context.Context, error) {
	return context.WithValue(ctx, startKey{}, time.Now()), nil
}

func (m *MetricsHook) AfterProcess(ctx context.Context, cmd redis.Cmder) error {
	if m.latencyMetricSender == nil {
		return nil
	}
	if start, ok := ctx.Value(startKey{}).(time.Time); ok {
		labels := map[string]string{
			monitoring.MethodLabel: cmd.Name(),
			monitoring.StatusLabel: GetStatus(cmd.Err()),
		}
		m.latencyMetricSender.Send(start, time.Now(), labels)

	}

	return nil
}

func (m *MetricsHook) BeforeProcessPipeline(ctx context.Context, cmds []redis.Cmder) (context.Context, error) {
	return context.WithValue(ctx, startKey{}, time.Now()), nil
}

func (m *MetricsHook) AfterProcessPipeline(ctx context.Context, cmds []redis.Cmder) error {
	// TODO: implement for pipelines
	return nil
}
func GetStatus(err error) string {
	switch err {
	case nil:
		return "ok"
	case redis.Nil:
		return "not_found"
	default:
		return "server_error"
	}
}
