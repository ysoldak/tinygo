//go:build cortexm
// +build cortexm

package main

import (
	"device/arm"
	"time"
)

// This example shows how to reset system on panic.
// In your application you may: reset, blink an LED to indicate failure, or do something else.

func main() {
	defer resetOnPanic()
	println("START")

	panicExplicit() // calls panic() explicitly
	// panicRuntime() // nil pointer access
}

func panicExplicit() {
	for i := 5; i >= 0; i-- {
		println(".")
		time.Sleep(time.Second)
	}
	panic("AAA!!!111")
}

func panicRuntime() {
	for i := 5; i >= 0; i-- {
		println(100/i)
		time.Sleep(time.Second)
	}
}

func resetOnPanic() {
	if r := recover(); r != nil {
		println("PANIC")
		time.Sleep(time.Second)
		arm.SystemReset()
	}
}
