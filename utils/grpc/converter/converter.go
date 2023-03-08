package converter

import (
	"time"

	"github.com/golang/protobuf/ptypes/wrappers"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type GRPCConverter struct{}

func NewGRPCConverter() *GRPCConverter {
	return &GRPCConverter{}
}

func (g *GRPCConverter) ToStringPointer(grpcStringPointer *wrappers.StringValue) *string {
	if grpcStringPointer == nil {
		return nil
	}
	str := new(string)
	*str = grpcStringPointer.GetValue()
	return str
}

func (g *GRPCConverter) ToGRPCStringPointer(stringPointer *string) *wrappers.StringValue {
	if stringPointer == nil {
		return nil
	}
	return &wrappers.StringValue{Value: *stringPointer}
}

func (g *GRPCConverter) ToInt32Pointer(grpcInt32Pointer *wrappers.Int32Value) *int32 {
	if grpcInt32Pointer == nil {
		return nil
	}
	integer := new(int32)
	*integer = grpcInt32Pointer.GetValue()
	return integer
}

func (g *GRPCConverter) ToInt64Pointer(grpcInt64Pointer *wrappers.Int64Value) *int64 {
	if grpcInt64Pointer == nil {
		return nil
	}
	integer := new(int64)
	*integer = grpcInt64Pointer.GetValue()
	return integer
}

func (g *GRPCConverter) ToIntPointerFromInt32Value(grpcInt32Pointer *wrappers.Int32Value) *int {
	if grpcInt32Pointer == nil {
		return nil
	}
	integer := new(int)
	*integer = int(grpcInt32Pointer.GetValue())
	return integer
}

func (g *GRPCConverter) ToEnumStringPointer(enum GRPCEnum) *string {
	if enum.Number() == 0 {
		return nil
	}
	str := new(string)
	*str = enum.String()
	return str
}

func (g *GRPCConverter) ToGRPCInt32Pointer(intPointer *int32) *wrappers.Int32Value {
	if intPointer == nil {
		return nil
	}
	return &wrappers.Int32Value{Value: *intPointer}
}

func (g *GRPCConverter) ToGRPCInt64Pointer(intPointer *int64) *wrappers.Int64Value {
	if intPointer == nil {
		return nil
	}
	return &wrappers.Int64Value{Value: *intPointer}
}

func (g *GRPCConverter) ToGRPCTimestamp(time time.Time) *timestamppb.Timestamp {
	timestampProto := timestamppb.New(time)
	return timestampProto
}

func (g *GRPCConverter) ToBoolWithDefault(grpcBoolPointer *wrappers.BoolValue, defaultValue bool) bool {
	if grpcBoolPointer == nil {
		return defaultValue
	}
	return grpcBoolPointer.GetValue()
}
