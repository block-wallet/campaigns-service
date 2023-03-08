package errors

import (
	"testing"

	"github.com/smartystreets/assertions"
	"github.com/smartystreets/goconvey/convey"
)

func TestNewClientError(t *testing.T) {
	convey.Convey("Given the parameters", t, func() {
		code := 200
		message := "error getting response"
		convey.Convey("When I ask for a new client error", func() {
			clientError := NewClient(code, message)
			convey.Convey("Then the message should have all the parameters", func() {
				convey.So(clientError, assertions.ShouldNotBeNil)
				convey.So(clientError.code, assertions.ShouldEqual, code)
				convey.So(clientError.message, assertions.ShouldEqual, message)
				convey.So(clientError.Error(), assertions.ShouldEqual, "HTTP client error with code: 200 and message: error getting response")
			})
		})
	})
}
