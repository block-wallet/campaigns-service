package validator

import (
	"strings"
	"testing"

	"google.golang.org/protobuf/types/known/structpb"

	"google.golang.org/grpc/status"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
)

func TestUuidForUuid(t *testing.T) {
	sv := SimpleValidation{
		Parameter: "b1dfe818-22fa-4940-a5ee-bb325b0a0e14",
		Validator: Uuid,
		ErrorMsg:  "ok",
	}

	firstError := Uuid(sv.Parameter)
	// Validation
	assert.NotNil(t, firstError)
	assert.Equal(t, firstError, true)
}

func TestUuidForNonUuid(t *testing.T) {
	cases := []SimpleValidation{{
		Parameter: nil,
		Validator: Uuid,
		ErrorMsg:  "nil",
	}, {
		Parameter: "",
		Validator: Uuid,
		ErrorMsg:  "empty",
	}, {
		Parameter: "other thing",
		Validator: Uuid,
		ErrorMsg:  "other thing",
	}}
	for _, c := range cases {
		c := c
		t.Run(c.ErrorMsg, func(t *testing.T) {
			t.Parallel()
			// Operation
			firstError := Uuid(c.Parameter)
			// Validation
			assert.NotNil(t, firstError)
			assert.Equal(t, firstError, false)
		})
	}
}

func TestZeroForZero(t *testing.T) {
	var int int
	var int8 int8
	var int16 int16
	var int32 int32
	var int64 int64
	var uint uint
	var uint8 uint8
	var uint16 uint16
	var uint32 uint32
	var uint64 uint64
	var float32 float32
	var float64 float64

	cases := []SimpleValidation{{
		Parameter: int,
		Validator: Zero,
		ErrorMsg:  "int 0",
	}, {
		Parameter: int8,
		Validator: Zero,
		ErrorMsg:  "int8 0",
	}, {
		Parameter: int16,
		Validator: Zero,
		ErrorMsg:  "int16 0",
	}, {
		Parameter: int32,
		Validator: Zero,
		ErrorMsg:  "int32 0",
	}, {
		Parameter: int64,
		Validator: Zero,
		ErrorMsg:  "int64 0",
	}, {
		Parameter: uint,
		Validator: Zero,
		ErrorMsg:  "uint 0",
	}, {
		Parameter: uint8,
		Validator: Zero,
		ErrorMsg:  "uint8 0",
	}, {
		Parameter: uint16,
		Validator: Zero,
		ErrorMsg:  "uint16 0",
	}, {
		Parameter: uint32,
		Validator: Zero,
		ErrorMsg:  "uint32 0",
	}, {
		Parameter: uint64,
		Validator: Zero,
		ErrorMsg:  "uint64 0",
	}, {
		Parameter: float32,
		Validator: Zero,
		ErrorMsg:  "float32 0",
	}, {
		Parameter: float64,
		Validator: Zero,
		ErrorMsg:  "float64 0",
	}}
	for _, c := range cases {
		c := c
		t.Run(c.ErrorMsg, func(t *testing.T) {
			t.Parallel()
			// Operation
			firstError := Zero(c.Parameter)
			// Validation
			assert.NotNil(t, firstError)
			assert.Equal(t, firstError, true)
		})
	}
}

func TestNZeroForNonZero(t *testing.T) {
	var int = 1
	var int8 int8 = 1
	var int16 int16 = 1
	var int32 int32 = 1
	var int64 int64 = 1
	var uint uint = 1
	var uint8 uint8 = 1
	var uint16 uint16 = 1
	var uint32 uint32 = 1
	var uint64 uint64 = 1
	var float32 float32 = 1
	var float64 float64 = 1

	cases := []SimpleValidation{{
		Parameter: int,
		Validator: Zero,
		ErrorMsg:  "int non 0",
	}, {
		Parameter: int8,
		Validator: Zero,
		ErrorMsg:  "int8 non 0",
	}, {
		Parameter: int16,
		Validator: Zero,
		ErrorMsg:  "int16 non 0",
	}, {
		Parameter: int32,
		Validator: Zero,
		ErrorMsg:  "int32 non 0",
	}, {
		Parameter: int64,
		Validator: Zero,
		ErrorMsg:  "int64 non 0",
	}, {
		Parameter: uint,
		Validator: Zero,
		ErrorMsg:  "uint non 0",
	}, {
		Parameter: uint8,
		Validator: Zero,
		ErrorMsg:  "uint8 non 0",
	}, {
		Parameter: uint16,
		Validator: Zero,
		ErrorMsg:  "uint16 non 0",
	}, {
		Parameter: uint32,
		Validator: Zero,
		ErrorMsg:  "uint32 non 0",
	}, {
		Parameter: uint64,
		Validator: Zero,
		ErrorMsg:  "uint64 non 0",
	}, {
		Parameter: float32,
		Validator: Zero,
		ErrorMsg:  "float32 non 0",
	}, {
		Parameter: float64,
		Validator: Zero,
		ErrorMsg:  "float64 non 0",
	}}
	for _, c := range cases {
		c := c
		t.Run(c.ErrorMsg, func(t *testing.T) {
			t.Parallel()
			// Operation
			firstError := Zero(c.Parameter)
			// Validation
			assert.NotNil(t, firstError)
			assert.Equal(t, firstError, false)
		})
	}
}

func TestNilForNil(t *testing.T) {
	var nilPointer *interface{} = nil
	var nilArray []string = nil
	var nilMap map[string]string = nil
	var nilStruct = structpb.NewNullValue()
	cases := []SimpleValidation{{
		Parameter: nil,
		Validator: Nil,
		ErrorMsg:  "nil",
	}, {
		Parameter: nilPointer,
		Validator: Nil,
		ErrorMsg:  "nil",
	}, {
		Parameter: nilArray,
		Validator: Nil,
		ErrorMsg:  "nil",
	}, {
		Parameter: nilMap,
		Validator: Nil,
		ErrorMsg:  "nil",
	}, {
		Parameter: nilStruct,
		Validator: Nil,
		ErrorMsg:  "nil",
	}}
	for _, c := range cases {
		c := c
		t.Run(c.ErrorMsg, func(t *testing.T) {
			t.Parallel()
			// Operation
			firstError := Nil(c.Parameter)
			// Validation
			assert.NotNil(t, firstError)
			assert.Equal(t, firstError, true)
		})
	}
}

func TestNilForNonNil(t *testing.T) {
	cases := []SimpleValidation{{
		Parameter: 1,
		Validator: Nil,
		ErrorMsg:  "non nil",
	}, {
		Parameter: 1.1,
		Validator: Nil,
		ErrorMsg:  "non nil",
	}, {
		Parameter: "1",
		Validator: Nil,
		ErrorMsg:  "non nil",
	}, {
		Parameter: true,
		Validator: Nil,
		ErrorMsg:  "non nil",
	}, {
		Parameter: []int{1, 2},
		Validator: Nil,
		ErrorMsg:  "non nil",
	}, {
		Parameter: map[string]string{"1": "1"},
		Validator: Nil,
		ErrorMsg:  "non nil",
	}, {
		Parameter: structpb.Value_StringValue{StringValue: "test"},
		Validator: Nil,
		ErrorMsg:  "non nil",
	}}
	for _, c := range cases {
		c := c
		t.Run(c.ErrorMsg, func(t *testing.T) {
			t.Parallel()
			// Operation
			firstError := Nil(c.Parameter)
			// Validation
			assert.NotNil(t, firstError)
			assert.Equal(t, firstError, false)
		})
	}
}

func TestNonZeroForZero(t *testing.T) {
	var int int
	var int8 int8
	var int16 int16
	var int32 int32
	var int64 int64
	var uint uint
	var uint8 uint8
	var uint16 uint16
	var uint32 uint32
	var uint64 uint64
	var float32 float32
	var float64 float64

	cases := []SimpleValidation{{
		Parameter: int,
		Validator: NonZero,
		ErrorMsg:  "int 0",
	}, {
		Parameter: int8,
		Validator: NonZero,
		ErrorMsg:  "int8 0",
	}, {
		Parameter: int16,
		Validator: NonZero,
		ErrorMsg:  "int16 0",
	}, {
		Parameter: int32,
		Validator: NonZero,
		ErrorMsg:  "int32 0",
	}, {
		Parameter: int64,
		Validator: NonZero,
		ErrorMsg:  "int64 0",
	}, {
		Parameter: uint,
		Validator: NonZero,
		ErrorMsg:  "uint 0",
	}, {
		Parameter: uint8,
		Validator: NonZero,
		ErrorMsg:  "uint8 0",
	}, {
		Parameter: uint16,
		Validator: NonZero,
		ErrorMsg:  "uint16 0",
	}, {
		Parameter: uint32,
		Validator: NonZero,
		ErrorMsg:  "uint32 0",
	}, {
		Parameter: uint64,
		Validator: NonZero,
		ErrorMsg:  "uint64 0",
	}, {
		Parameter: float32,
		Validator: NonZero,
		ErrorMsg:  "float32 0",
	}, {
		Parameter: float64,
		Validator: NonZero,
		ErrorMsg:  "float64 0",
	}}
	for _, c := range cases {
		c := c
		t.Run(c.ErrorMsg, func(t *testing.T) {
			t.Parallel()
			// Operation
			firstError := NonZero(c.Parameter)
			// Validation
			assert.NotNil(t, firstError)
			assert.Equal(t, firstError, false)
		})
	}
}

func TestNonZeroForNonZero(t *testing.T) {
	var int = 1
	var int8 int8 = 1
	var int16 int16 = 1
	var int32 int32 = 1
	var int64 int64 = 1
	var uint uint = 1
	var uint8 uint8 = 1
	var uint16 uint16 = 1
	var uint32 uint32 = 1
	var uint64 uint64 = 1
	var float32 float32 = 1
	var float64 float64 = 1

	cases := []SimpleValidation{{
		Parameter: int,
		Validator: NonZero,
		ErrorMsg:  "int non 0",
	}, {
		Parameter: int8,
		Validator: NonZero,
		ErrorMsg:  "int8 non 0",
	}, {
		Parameter: int16,
		Validator: NonZero,
		ErrorMsg:  "int16 non 0",
	}, {
		Parameter: int32,
		Validator: NonZero,
		ErrorMsg:  "int32 non 0",
	}, {
		Parameter: int64,
		Validator: NonZero,
		ErrorMsg:  "int64 non 0",
	}, {
		Parameter: uint,
		Validator: NonZero,
		ErrorMsg:  "uint non 0",
	}, {
		Parameter: uint8,
		Validator: NonZero,
		ErrorMsg:  "uint8 non 0",
	}, {
		Parameter: uint16,
		Validator: NonZero,
		ErrorMsg:  "uint16 non 0",
	}, {
		Parameter: uint32,
		Validator: NonZero,
		ErrorMsg:  "uint32 non 0",
	}, {
		Parameter: uint64,
		Validator: NonZero,
		ErrorMsg:  "uint64 non 0",
	}, {
		Parameter: float32,
		Validator: NonZero,
		ErrorMsg:  "float32 non 0",
	}, {
		Parameter: float64,
		Validator: NonZero,
		ErrorMsg:  "float64 non 0",
	}}
	for _, c := range cases {
		c := c
		t.Run(c.ErrorMsg, func(t *testing.T) {
			t.Parallel()
			// Operation
			firstError := NonZero(c.Parameter)
			// Validation
			assert.NotNil(t, firstError)
			assert.Equal(t, firstError, true)
		})
	}
}

func TestArrayNoNEmptyForEmptyArray(t *testing.T) {
	var emptyArray []string

	sv := SimpleValidation{
		Parameter: emptyArray,
		Validator: ArrayNoNEmpty,
		ErrorMsg:  "doesn't pass",
	}

	firstError := FirstNonValid(sv)

	assert.NotNil(t, firstError)
}

func TestArrayNoNEmptyForNonEmptyArray(t *testing.T) {
	nonEmptyArray := []string{"A", "B"}

	sv := SimpleValidation{
		Parameter: nonEmptyArray,
		Validator: ArrayNoNEmpty,
		ErrorMsg:  "doesn't pass",
	}

	firstError := FirstNonValid(sv)

	assert.Nil(t, firstError)
}

func TestEmptyIsInvalidForStringPresent(t *testing.T) {
	var emptyThing string

	sv := SimpleValidation{
		Parameter: emptyThing,
		Validator: StringPresent,
		ErrorMsg:  "doesn't pass",
	}

	firstError := FirstNonValid(sv)

	assert.NotNil(t, firstError)
}

func TestNonEmptyIsInvalidForStringPresent(t *testing.T) {
	var fullThing string = "I'm full of value"
	sv := SimpleValidation{
		Parameter: fullThing,
		Validator: StringPresent,
		ErrorMsg:  "this shouldn't happen",
	}

	firstError := FirstNonValid(sv)

	assert.Nil(t, firstError)
}

func TestNonNilValidatorValidatesNonNilParam(t *testing.T) {
	pass := SimpleValidation{
		Parameter: "something",
		Validator: NonNil,
		ErrorMsg:  "this shouldn't happen",
	}

	firstError := FirstNonValid(pass)

	assert.Nil(t, firstError)

	breaks := SimpleValidation{
		Parameter: nil,
		Validator: NonNil,
		ErrorMsg:  "this should pop",
	}

	firstError = FirstNonValid(breaks)

	assert.NotNil(t, firstError)

	var nilPointer *interface{} = nil

	breaks = SimpleValidation{
		Parameter: nilPointer,
		Validator: NonNil,
		ErrorMsg:  "this should pop",
	}

	firstError = FirstNonValid(breaks)

	assert.NotNil(t, firstError)
}

func TestErrorsAreHonored(t *testing.T) {
	theMessage := "doesn't pass"
	var emptyThing string

	sv := SimpleValidation{
		Parameter: emptyThing,
		Validator: StringPresent,
		ErrorMsg:  theMessage,
	}

	err := FirstNonValid(sv)

	assert.Equal(t, theMessage, err.Error())
	assert.Equal(t, codes.InvalidArgument, status.Code(err.ToGRPCError()))
}

func TestManyValidParametersPass(t *testing.T) {
	fullThing1 := "I'm full of value"
	fullThing2 := "I'm full of value too"

	sv1 := SimpleValidation{
		Parameter: fullThing1,
		Validator: StringPresent,
		ErrorMsg:  "this shouldn't happen",
	}

	sv2 := SimpleValidation{
		Parameter: fullThing2,
		Validator: StringPresent,
		ErrorMsg:  "this shouldn't happen",
	}

	fnv := FirstNonValid(sv1, sv2)

	assert.Nil(t, fnv)
}

func TestInvalidParametersOnSameValidator(t *testing.T) {
	var emptyThing string

	sv1 := SimpleValidation{
		Parameter: emptyThing,
		ErrorMsg:  "this should arise",
	}

	var fullThing string = "I'm full of value"
	sv2 := SimpleValidation{
		Parameter: fullThing,
		ErrorMsg:  "this shouldn't happen",
	}

	fnv := FirstNonValidSameValidation(StringPresent, sv1, sv2)
	assert.NotNil(t, fnv)

	assert.Equal(t, "this should arise", fnv.Error())
	assert.Equal(t, codes.InvalidArgument, status.Code(fnv.ToGRPCError()))
}

func TestManyValidParametersPassOnSameValidator(t *testing.T) {
	fullThing1 := "I'm full of value"
	fullThing2 := "I'm full of value too"

	sv1 := SimpleValidation{
		Parameter: fullThing1,
		ErrorMsg:  "this shouldn't happen",
	}

	sv2 := SimpleValidation{
		Parameter: fullThing2,
		ErrorMsg:  "this shouldn't happen",
	}

	fnv := FirstNonValidSameValidation(StringPresent, sv1, sv2)

	assert.Nil(t, fnv)
}

type MyLowerCaseEquality struct {
	compareTo string
}

func (mle *MyLowerCaseEquality) lowerCaseEquality(that interface{}) bool {
	return strings.Compare(mle.compareTo, strings.ToLower(that.(string))) == 0
}

func TestCustomValidator(t *testing.T) {
	something := "I'm something"

	mle := MyLowerCaseEquality{}
	mle.compareTo = "i'm something"

	sv := SimpleValidation{
		Parameter: something,
		Validator: mle.lowerCaseEquality,
		ErrorMsg:  "this shouldn't happen",
	}

	firstError := FirstNonValid(sv)

	assert.Nil(t, firstError)
}

func TestEqualityValidator(t *testing.T) {
	eqc := EqualityComparison{
		CompareTo: 10,
	}

	otherThing := 10

	sv := SimpleValidation{
		Parameter: otherThing,
		Validator: eqc.Equal,
		ErrorMsg:  "this shouldn't happen",
	}

	firstError := FirstNonValid(sv)

	assert.Nil(t, firstError)
}

func TestNonEqualValidator(t *testing.T) {
	eqc := EqualityComparison{
		CompareTo: 10,
	}

	otherThing := 20

	sv := SimpleValidation{
		Parameter: otherThing,
		Validator: eqc.NonEqual,
		ErrorMsg:  "this shouldn't happen",
	}

	firstError := FirstNonValid(sv)

	assert.Nil(t, firstError)
}

func TestDeepNonEqualityValidator(t *testing.T) {
	eqc := EqualityComparison{
		CompareTo: map[string]string{
			"k1": "v1",
			"k2": "v2",
		},
	}

	otherThing := map[string]string{
		"k3": "v1",
		"k4": "v2",
	}

	sv := SimpleValidation{
		Parameter: otherThing,
		Validator: eqc.DeepNonEqual,
		ErrorMsg:  "this shouldn't happen",
	}

	firstError := FirstNonValid(sv)

	assert.Nil(t, firstError)
}

func TestDeepEqualityValidator(t *testing.T) {
	eqc := EqualityComparison{
		CompareTo: map[string]string{
			"k1": "v1",
			"k2": "v2",
		},
	}

	otherThing := map[string]string{
		"k1": "v1",
		"k2": "v2",
	}

	sv := SimpleValidation{
		Parameter: otherThing,
		Validator: eqc.DeepEqual,
		ErrorMsg:  "this shouldn't happen",
	}

	firstError := FirstNonValid(sv)

	assert.Nil(t, firstError)
}

func TestParameterNotInSet(t *testing.T) {
	var notInSet string = "I'm not in the following set"

	theSet := InSet{
		ThingInSet: []EqualityComparison{
			{
				CompareTo: "rewarded_video",
			},
			{
				CompareTo: "interstitial",
			},
			{
				CompareTo: "banner",
			},
		},
	}

	sv := SimpleValidation{
		Parameter: notInSet,
		Validator: theSet.InSet,
		ErrorMsg:  "doesn't pass",
	}

	firstError := FirstNonValid(sv)

	assert.NotNil(t, firstError)
}

func TestParameterInSet(t *testing.T) {
	var notInSet string = "rewarded_video"

	theSet := InSet{
		ThingInSet: []EqualityComparison{
			{
				CompareTo: "rewarded_video",
			},
			{
				CompareTo: "interstitial",
			},
			{
				CompareTo: "banner",
			},
		},
	}

	sv := SimpleValidation{
		Parameter: notInSet,
		Validator: theSet.InSet,
		ErrorMsg:  "this shouldn't pop up",
	}

	firstError := FirstNonValid(sv)

	assert.Nil(t, firstError)
}
