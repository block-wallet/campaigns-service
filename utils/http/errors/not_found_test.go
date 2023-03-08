package errors

import (
	"testing"

	"github.com/smartystreets/assertions"
	"github.com/smartystreets/goconvey/convey"
)

func TestNewNotFoundError(t *testing.T) {
	convey.Convey("Given the parameters", t, func() {
		message := "entity not found"
		convey.Convey("When I ask for a new not found error", func() {
			notFoundError := NewNotFound(message)
			convey.Convey("Then the message should have all the parameters", func() {
				convey.So(notFoundError, assertions.ShouldNotBeNil)
				convey.So(notFoundError.message, assertions.ShouldEqual, message)
				convey.So(notFoundError.Error(), assertions.ShouldEqual, "HTTP not found error with message: entity not found")
			})
		})
	})
}
