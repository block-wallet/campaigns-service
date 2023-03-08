package counter

type MetricSender interface {
	Send(labels map[string]string)
}
