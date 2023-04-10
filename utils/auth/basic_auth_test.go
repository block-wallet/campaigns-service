package auth

import (
	"context"
	"encoding/base64"
	"fmt"
	"testing"

	"github.com/smartystreets/goconvey/convey"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func newContext(username, password string) context.Context {
	credentials := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%v:%v", username, password)))
	return metadata.NewIncomingContext(context.Background(), metadata.New(map[string]string{"Authorization": fmt.Sprintf("Basic %v", credentials)}))
}

func Test_BasicAuth(t *testing.T) {
	username, password := "user1", "password1"
	basicAuth := NewBasicAuth(username, password)
	convey.Convey("Given an incoming request with auth header", t, func() {
		convey.Convey("When the credentials are not valid", func() {
			ctx := newContext("fakeUser", "fakePassword")
			convey.Convey("Then the validation should fail", func() {
				err := basicAuth.AuthenticateUsingContext(ctx)
				convey.So(err, convey.ShouldNotEqual, nil)
				convey.So(status.Code(err.ToGRPCError()), convey.ShouldEqual, codes.Unauthenticated)
			})
		})
		convey.Convey("When the credentials are valid", func() {
			ctx := newContext(username, password)
			convey.Convey("Then the validation should pass", func() {
				err := basicAuth.AuthenticateUsingContext(ctx)
				convey.So(err, convey.ShouldEqual, nil)
			})
		})
	})
	convey.Convey("Given an incoming without auth header", t, func() {
		convey.Convey("When the context has not authorization header", func() {
			convey.Convey("Then the validation should fail", func() {
				err := basicAuth.AuthenticateUsingContext(context.Background())
				convey.So(err, convey.ShouldNotEqual, nil)
			})
		})
	})
}
