package validator

import (
	"reflect"
	"regexp"

	"google.golang.org/protobuf/types/known/structpb"

	"github.com/block-wallet/golang-service-template/utils/errors"
)

type Validator func(interface{}) bool

type SimpleValidation struct {
	Parameter interface{}
	Validator Validator
	ErrorMsg  string
}

func (v SimpleValidation) Err() errors.RichError {
	return errors.NewInvalidArgument(v.ErrorMsg)
}

type Func struct {
	Fn func(interface{}) bool
}

func (f *Func) Satisfy(that interface{}) bool {
	return f.Fn(that)
}

type InSet struct {
	ThingInSet []EqualityComparison
}

func (i *InSet) InSet(that interface{}) bool {
	for _, e := range i.ThingInSet {
		if e.Equal(that) {
			return true
		}
	}
	return false
}

type EqualityComparison struct {
	CompareTo interface{}
}

func (ec *EqualityComparison) NonEqual(that interface{}) bool {
	return !ec.Equal(that)
}

func (ec *EqualityComparison) DeepNonEqual(that interface{}) bool {
	return !ec.DeepEqual(that)
}

func (ec *EqualityComparison) Equal(that interface{}) bool {
	return that == ec.CompareTo
}

func (ec *EqualityComparison) DeepEqual(that interface{}) bool {
	return reflect.DeepEqual(that, ec.CompareTo)
}

func Nil(param interface{}) bool {
	if param == nil {
		return true
	}

	switch param.(type) {
	case *structpb.Value:
		v := param.(*structpb.Value)
		if _, ok := v.GetKind().(*structpb.Value_NullValue); ok {
			return true
		}
	}

	switch reflect.TypeOf(param).Kind() {
	case reflect.Ptr, reflect.Map, reflect.Array, reflect.Chan, reflect.Slice:
		return reflect.ValueOf(param).IsNil()
	}

	return false
}

func NonNil(param interface{}) bool {
	return !Nil(param)
}

func NonZero(param interface{}) bool {
	return !Zero(param)
}

func Zero(param interface{}) bool {
	if param == nil {
		return false
	}

	switch reflect.TypeOf(param).Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Float32, reflect.Float64:
		return reflect.ValueOf(param).IsZero()
	}

	return false
}

func Uuid(param interface{}) bool {
	if param == nil {
		return false
	}
	if len(param.(string)) <= 0 {
		return false
	}

	r := regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$")
	return r.MatchString(param.(string))
}

func ArrayNoNEmpty(param interface{}) bool {
	if param == nil {
		return false
	}

	switch reflect.TypeOf(param).Kind() {
	case reflect.Array, reflect.Slice, reflect.Map, reflect.Chan:
		return reflect.ValueOf(param).Len() != 0
	}

	return true
}

func StringPresent(param interface{}) bool {
	return len(param.(string)) > 0
}

func ValidateParameter(param interface{}, validator func(interface{}) bool) bool {
	return validator(param)
}

func FirstNonValid(validations ...SimpleValidation) errors.RichError {
	for _, v := range validations {
		if !ValidateParameter(v.Parameter, v.Validator) {
			return v.Err()
		}
	}
	return nil
}

func FirstNonValidSameValidation(singleValidator Validator, validations ...SimpleValidation) errors.RichError {
	for _, v := range validations {
		if !ValidateParameter(v.Parameter, singleValidator) {
			return v.Err()
		}
	}
	return nil
}
