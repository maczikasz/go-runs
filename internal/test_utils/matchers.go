package test_utils

import "reflect"

func ShouldMatch(actual interface{}, expected ...interface{}) string {
	if len(expected) != 1 {
		return "This assertion requires exactly 1 expected value"
	}

	v := reflect.ValueOf(actual)
	if v.Kind() != reflect.Slice && v.Kind() != reflect.Array {
		return "This assertion requires a slice or an array"
	}

	f := reflect.ValueOf(expected[0])
	if f.Kind() != reflect.Func {
		return "This assertion requires a matcher function"
	}
	for i := 0; i < v.Len(); i++ {
		elem := v.Index(i)
		matches := f.Call([]reflect.Value{elem})
		if len(matches) != 1 || matches[0].Kind() != reflect.Bool {
			return "The assert function must return exactly 1 boolean"
		}
		if matches[0].Bool() {
			return ""
		}
	}

	return "The slice did not contain a matching object"

}
