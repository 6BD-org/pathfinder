package utils

import "reflect"

// Filter out elements using predicates
func Filter(l interface{}, predicate func(l interface{}) bool, t reflect.Type) []interface{} {
	res := make([]interface{}, 0)
	switch reflect.TypeOf(l).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(l)
		for i := 0; i < s.Len(); i++ {
			if predicate(s.Index(i).Convert(t).Interface()) {
				res = append(res, s.Index(i).Interface())
			}
		}
	}
	return res
}
