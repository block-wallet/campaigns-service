package errors

import (
	"testing"

	"github.com/smartystreets/assertions"
	"github.com/smartystreets/goconvey/convey"
)

func TestNewServerError(t *testing.T) {
	convey.Convey("Given the parameters", t, func() {
		code := 200
		message := "error getting response"
		convey.Convey("When I ask for a new server error", func() {
			serverError := NewServer(code, message)
			convey.Convey("Then the message should have all the parameters", func() {
				convey.So(serverError, assertions.ShouldNotBeNil)
				convey.So(serverError.code, assertions.ShouldEqual, code)
				convey.So(serverError.message, assertions.ShouldEqual, message)
				convey.So(serverError.Error(), assertions.ShouldEqual, "HTTP server error with code: 200 and message: error getting response")
			})
		})
	})
}
