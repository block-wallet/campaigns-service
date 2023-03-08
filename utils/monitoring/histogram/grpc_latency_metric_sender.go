package histogram

import (
	"time"

	"google.golang.org/grpc/status"
)

type GRPCLatencyMetricSender struct {
	latencyMetricSender LatencyMetricSender
}

func NewGRPCLatencyMetricSender(latencyMetricSender LatencyMetricSender) *GRPCLatencyMetricSender {
	return &GRPCLatencyMetricSender{
		latencyMetricSender: latencyMetricSender,
	}
}

func (g *GRPCLatencyMetricSender) Send(start, end time.Time, method string, req interface{}, err error) {
	labels := g.getLabelsByMethod(method, status.Code(err).String(), req)
	g.latencyMetricSender.Send(start, end, labels)
}

func (g *GRPCLatencyMetricSender) getLabelsByMethod(method, status string, req interface{}) map[string]string {
	labels := make(map[string]string)
	labels["method"] = method
	labels["status"] = status
	return labels
}
