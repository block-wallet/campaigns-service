package histogram

import (
	"testing"

	"github.com/smartystreets/assertions"
	"github.com/smartystreets/goconvey/convey"
)

func Test_GetLabelsByMethod(t *testing.T) {
	convey.Convey("Given an incoming request with OK response", t, func() {

		grpcLatencyMetricSender := NewGRPCLatencyMetricSender(nil)
		convey.Convey("When we get the labels", func() {
			labels := grpcLatencyMetricSender.getLabelsByMethod("/health", "ok", nil)
			convey.Convey("Then we should have all labels set", func() {
				convey.So(len(labels), assertions.ShouldEqual, 2)
				convey.So(labels["method"], assertions.ShouldEqual, "/health")
				convey.So(labels["status"], assertions.ShouldEqual, "ok")
			})
		})
	})
}
