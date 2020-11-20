package main

import (
	"log"
	"testing"
)

func Foo() (string, error) {
	return "", nil
}

type B struct {
	Val int
}

type A struct {
	B B
}

func TestVariable(t *testing.T) {
	b := B{Val: 1}
	a := A{B: b}
	a.B.Val = 3
	log.Println(a)
}
