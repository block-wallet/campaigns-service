package healthgrpcservice_test

import (
	"context"
	"testing"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/smartystreets/goconvey/convey"

	healthgrpcservice "github.com/block-wallet/golang-service-template/domain/health/service/grpc"
	ethservicev1health "github.com/block-wallet/golang-service-template/protos/ethservicev1/src/health"
	. "github.com/smartystreets/assertions"
)

func TestHealthHandler(t *testing.T) {
	convey.Convey("Given a health service handler", t, func() {
		healthService := healthgrpcservice.NewHandler()
		convey.Convey("When i ask for the status", func() {
			responseMsg, err := healthService.Status(context.Background(), new(empty.Empty))
			convey.Convey("Then the status is ALIVE", func() {
				convey.So(err, ShouldBeNil)
				convey.So(responseMsg, ShouldNotBeNil)
				convey.So(responseMsg.Status, ShouldEqual, ethservicev1health.HealthStatus_HEALTH_STATUS_ALIVE)
			})
		})
	})
}
