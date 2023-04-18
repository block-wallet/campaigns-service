package healthgrpcservice_test

import (
	"context"
	"testing"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/smartystreets/goconvey/convey"

	healthgrpcservice "github.com/block-wallet/campaigns-service/domain/health/service/grpc"
	campaignsservicev1health "github.com/block-wallet/campaigns-service/protos/src/campaignsservicev1/health"
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
				convey.So(responseMsg.Status, ShouldEqual, campaignsservicev1health.HealthStatus_HEALTH_STATUS_ALIVE)
			})
		})
	})
}
