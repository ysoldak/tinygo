package main

/*
#define foo 3
#define bar foo
#define unreferenced 4
#define referenced unreferenced
*/
import "C"

const (
	Foo = C.foo
	Bar = C.bar
	Baz = C.referenced
)
