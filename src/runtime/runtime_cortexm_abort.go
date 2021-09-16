//go:build cortexm && !nxp && !qemu
// +build cortexm,!nxp,!qemu

package runtime

import (
	"device/arm"
)

func exit(code int) {
	abort()
}

var ResetOnAbort = false

func abort() {
	if ResetOnAbort {
		arm.SystemReset()
	}
	// lock up forever
	for {
		arm.Asm("wfi")
	}
}
