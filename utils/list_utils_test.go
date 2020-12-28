package utils

import (
	"log"
	"reflect"
	"testing"
)

func TestFilter(t *testing.T) {
	var source = []int{1, 2, 3}
	res := Filter(source, func(l interface{}) bool { return l.(int) > 1 }, reflect.TypeOf(1))
	for _, v := range res {
		log.Println(v.(int))
	}
}
