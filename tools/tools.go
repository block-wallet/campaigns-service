//go:build tools
// +build tools

package tools

import (
	_ "github.com/golang/protobuf/protoc-gen-go"
	_ "github.com/golangci/golangci-lint/cmd/golangci-lint"
	_ "github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway"
	_ "github.com/pseudomuto/protoc-gen-doc"
	_ "github.com/pseudomuto/protoc-gen-doc/extensions/envoyproxy_validate"
	_ "github.com/pseudomuto/protoc-gen-doc/extensions/validator_field"
	_ "golang.org/x/tools/cmd/goimports"
	_ "mvdan.cc/gofumpt/format"
)
