package main

import "testing"

func Foo() (string, error) {
	return "", nil
}

func TestVariable(t *testing.T) {
	a, err := Foo()

	b, err := Foo()
	a, err := Foo()
}
