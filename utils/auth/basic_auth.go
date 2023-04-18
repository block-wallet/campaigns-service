package auth

import (
	"context"
	"encoding/base64"
	"strings"

	"github.com/block-wallet/campaigns-service/utils/errors"
	"google.golang.org/grpc/metadata"
)

type BasicAuth struct {
	username string
	password string
}

func NewBasicAuth(username string, password string) *BasicAuth {
	return &BasicAuth{
		username: username,
		password: password,
	}
}

func (a *BasicAuth) AuthenticateUsingContext(ctx context.Context) errors.RichError {
	md, _ := metadata.FromIncomingContext(ctx)
	authroization := md.Get("Authorization")
	authenticated := false
	if len(authroization) > 0 {
		authHeader := authroization[0]
		if splittedHeader := strings.Split(authHeader, " "); len(splittedHeader) == 2 {
			if splittedHeader[0] == "Basic" {
				decodedCredentials, err := base64.StdEncoding.DecodeString(splittedHeader[1])
				if err == nil {
					cred := strings.Split(string(decodedCredentials), ":")
					if len(cred) == 2 {
						username := cred[0]
						password := cred[1]
						if username == a.username && password == a.password {
							authenticated = true
						}
					}

				}
			}
		}

	}

	if !authenticated {
		return errors.NewUnauthenticated("req not authenticated")
	}
	return nil
}
