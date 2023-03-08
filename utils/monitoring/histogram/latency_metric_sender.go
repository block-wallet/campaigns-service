package histogram

import (
	"time"
)

type LatencyMetricSender interface {
	Send(start, end time.Time, labels map[string]string)
}
