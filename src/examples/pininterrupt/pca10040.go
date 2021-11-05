//go:build pca10040
// +build pca10040

package main

import "machine"

const (
	buttonPin       = machine.BUTTON
	buttonMode      = machine.PinInputPullup
	buttonPinChange = machine.PinRising
)
