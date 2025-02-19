package functools

import (
	"errors"
	"reflect"
)

// Reduce applys a function(the first parameter) of two arguments cumulatively to each element of a slice(second parameter), and the initial value is the third parameter.
func Reduce(function, slice, initial interface{}) (ret interface{}, err error) {
	err = nil
	defer func() {
		if interfaceErr := recover(); interfaceErr != nil {
			err = errors.New(interfaceErr.(string))
		}
	}()
	ret = reduce(function, slice, initial)
	return
}

func reduce(function, slice, initial interface{}) interface{} {
	in := reflect.ValueOf(slice)
	if in.Kind() != reflect.Slice {
		panic("reduce: The first param is not a slice")
	}
	fn := reflect.ValueOf(function)
	inType := in.Type().Elem()
	if inType != reflect.TypeOf(initial) {
		panic("reduce: The Type of first param and elements in second param must be the same")
	}
	if !verifyReduceFuncType(fn, inType) {
		panic("reduce: Function must be of type func(" + inType.String() + ")" + inType.String())
	}
	var param [2]reflect.Value
	out := reflect.ValueOf(initial)
	for i := 0; i < in.Len(); i++ {
		param[0] = out
		param[1] = in.Index(i)
		out = fn.Call(param[:])[0]
	}
	return out.Interface()
}

func verifyReduceFuncType(fn reflect.Value, elemType reflect.Type) bool {
	if fn.Kind() != reflect.Func {
		return false
	}
	if fn.Type().NumIn() != 2 || fn.Type().NumOut() != 1 {
		return false
	}
	if elemType != fn.Type().In(0) || fn.Type().In(0) != fn.Type().In(1) || fn.Type().In(1) != fn.Type().Out(0) {
		return false
	}
	return true
}
