package converter

import "google.golang.org/protobuf/reflect/protoreflect"

type GRPCEnum interface {
	String() string
	Number() protoreflect.EnumNumber
}
