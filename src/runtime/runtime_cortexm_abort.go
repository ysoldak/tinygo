// +build cortexm,!nxp,!qemu

package runtime

import (
	"github.com/sago35/device/arm"
)

func abort() {
	// disable all interrupts
	arm.DisableInterrupts()

	// lock up forever
	for {
		arm.Asm("wfi")
	}
}
