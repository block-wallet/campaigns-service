package histogram

import (
	"time"
)

type RequestLatencyMetricSender interface {
	Send(start, end time.Time, method string, req interface{}, err error)
}
