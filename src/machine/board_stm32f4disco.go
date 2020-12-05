// +build stm32f4disco

package machine

import (
	"device/stm32"
	"runtime/interrupt"
)

const (
	LED         = LED_BUILTIN
	LED1        = LED_GREEN
	LED2        = LED_ORANGE
	LED3        = LED_RED
	LED4        = LED_BLUE
	LED_BUILTIN = LED_GREEN
	LED_GREEN   = PD12
	LED_ORANGE  = PD13
	LED_RED     = PD14
	LED_BLUE    = PD15
)

// UART pins
const (
	UART_TX_PIN = PA2
	UART_RX_PIN = PA3
)

var (
	UART0 = UART{
		Buffer:          NewRingBuffer(),
		Bus:             stm32.USART2,
		AltFuncSelector: stm32.AF7_USART1_2_3,
	}
	UART1 = &UART0
)

// set up RX IRQ handler. Follow similar pattern for other UARTx instances
func init() {
	UART0.Interrupt = interrupt.New(stm32.IRQ_USART2, UART0.handleInterrupt)
}

// SPI pins
const (
	SPI1_SCK_PIN = PA5
	SPI1_SDI_PIN = PA6
	SPI1_SDO_PIN = PA7
	SPI0_SCK_PIN = SPI1_SCK_PIN
	SPI0_SDI_PIN = SPI1_SDI_PIN
	SPI0_SDO_PIN = SPI1_SDO_PIN
)

// MEMs accelerometer
const (
	MEMS_ACCEL_CS   = PE3
	MEMS_ACCEL_INT1 = PE0
	MEMS_ACCEL_INT2 = PE1
)

// Since the first interface is named SPI1, both SPI0 and SPI1 refer to SPI1.
// TODO: implement SPI2 and SPI3.
var (
	SPI0 = SPI{
		Bus:             stm32.SPI1,
		AltFuncSelector: stm32.AF5_SPI1_SPI2,
	}
	SPI1 = &SPI0
)

// -- I2C ----------------------------------------------------------------------

const (
	// #===========#==========#==============#==============#=======#=======#
	// | Interface | Hardware |  Bus(Freq)   | SDA/SCL Pins | AltFn | Alias |
	// #===========#==========#==============#==============#=======#=======#
	// |   I2C1    |   I2C1   | APB1(42 MHz) |   PB7/PB6    |   4   |   ~   |
	// |   I2C2    |   I2C2   | APB1(42 MHz) |  PB11/PB10   |   4   |   ~   |
	// |   I2C3    |   I2C1   | APB1(42 MHz) |   PB7/PB6    |   4   |   ~   |
	// | --------- | -------- | ------------ | ------------ | ----- | ----- |
	// |   I2C0    |   I2C1   | APB1(42 MHz) |   PB7/PB6    |   4   | I2C1  |
	// #===========#==========#==============#==============#=======#=======#
	NUM_I2C_INTERFACES = 3

	I2C1_SDA_PIN = PB7 // I2C1 = hardware: I2C1
	I2C1_SCL_PIN = PB6 //

	I2C2_SDA_PIN = PB11 // I2C2 = hardware: I2C2
	I2C2_SCL_PIN = PB10 //

	I2C3_SDA_PIN = PB7 // I2C3 = hardware: I2C1
	I2C3_SCL_PIN = PB6 //  (interface duplicated on second pair of pins)

	I2C0_SDA_PIN = I2C1_SDA_PIN // I2C0 = alias: I2C1
	I2C0_SCL_PIN = I2C1_SCL_PIN //

	I2C_SDA_PIN = I2C0_SDA_PIN // default/primary I2C pins
	I2C_SCL_PIN = I2C0_SCL_PIN //
)

var (
	I2C1 = I2C{
		Bus:             stm32.I2C1,
		AltFuncSelector: stm32.AF4_I2C1_2_3,
	}
	I2C2 = I2C{
		Bus:             stm32.I2C2,
		AltFuncSelector: stm32.AF4_I2C1_2_3,
	}
	I2C3 = I2C{
		Bus:             stm32.I2C1,
		AltFuncSelector: stm32.AF4_I2C1_2_3,
	}
	I2C0 = I2C1
)
