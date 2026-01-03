package lcd

import (
	"fmt"
	"os"
	"strings"
	"time"

	"periph.io/x/conn/v3/i2c"
	"periph.io/x/conn/v3/i2c/i2creg"
	"periph.io/x/host/v3"
)

const (
	// I2C LCD address (can be detected via i2cdetect).
	lcdAddr = 0x27

	// I2C bus 4 - software I2C on GPIO 23 (SDA) and GPIO 24 (SCL).
	// Configured in /boot/config.txt: dtoverlay=i2c-gpio,bus=4,i2c_gpio_sda=23,i2c_gpio_scl=24
	i2cBus = "4"

	// PCF8574 pin mapping for I2C LCD backpack.
	// P0=RS, P1=RW, P2=E, P3=Backlight, P4-P7=D4-D7
	pinRS        = 0x01
	pinRW        = 0x02
	pinE         = 0x04
	pinBacklight = 0x08

	// HD44780 commands.
	cmdClear        = 0x01
	cmdHome         = 0x02
	cmdEntryMode    = 0x04
	cmdDisplayCtrl  = 0x08
	cmdFunctionSet  = 0x20
	cmdSetDDRAMAddr = 0x80

	// Flags for display entry mode.
	entryLeft = 0x02

	// Flags for display on/off control.
	displayOn = 0x04
	cursorOff = 0x00
	blinkOff  = 0x00

	// Flags for function set.
	mode4Bit = 0x00
	lines2   = 0x08
	dots5x8  = 0x00

	// Display geometry.
	Cols = 16
	Rows = 2
)

var (
	i2cDev    *i2c.Dev
	backlight byte = pinBacklight // Backlight state
)

// Stop cleans up the LCD device
func Stop() {
	if i2cDev != nil {
		Clear()
	}
}

// Start initializes the I2C LCD.
func Start() {
	// Initialize periph.io host drivers.
	fmt.Println("Initializing periph.io...")
	if _, err := host.Init(); err != nil {
		panic(err)
	}

	// Open I2C bus 4 (software I2C on GPIO 23/24).
	fmt.Printf("Opening I2C bus %s...\n", i2cBus)
	bus, err := i2creg.Open(i2cBus)
	if err != nil {
		panic(fmt.Sprintf("failed to open I2C bus: %v", err))
	}

	// Create device handle for LCD at its address.
	fmt.Printf("Connecting to I2C LCD at address 0x%02x...\n", lcdAddr)
	i2cDev = &i2c.Dev{Bus: bus, Addr: uint16(lcdAddr)}

	// Initialize LCD in 4-bit mode.
	fmt.Println("Initializing LCD in 4-bit mode...")

	// Wait for LCD to power up.
	time.Sleep(50 * time.Millisecond)

	// Initialization sequence for 4-bit mode.
	// See HD44780 datasheet page 46.
	write4bits(0x03 << 4) // Function set: 8-bit mode
	time.Sleep(5 * time.Millisecond)

	write4bits(0x03 << 4) // Function set: 8-bit mode
	time.Sleep(150 * time.Microsecond)

	write4bits(0x03 << 4) // Function set: 8-bit mode
	time.Sleep(150 * time.Microsecond)

	write4bits(0x02 << 4) // Function set: 4-bit mode
	time.Sleep(150 * time.Microsecond)

	// Now in 4-bit mode, configure the display.
	writeCommand(cmdFunctionSet | mode4Bit | lines2 | dots5x8)      // 4-bit, 2 lines, 5x8 font
	writeCommand(cmdDisplayCtrl | displayOn | cursorOff | blinkOff) // Display on, cursor off, blink off
	writeCommand(cmdClear)                                          // Clear display
	time.Sleep(2 * time.Millisecond)                                // Clear command takes longer
	writeCommand(cmdEntryMode | entryLeft)                          // Entry mode: left to right

	fmt.Println("LCD initialization complete")
}

// write4bits sends 4 bits to the LCD via I2C.
func write4bits(data byte) {
	// Combine data with current backlight state.
	value := data | backlight

	// Send data with E high.
	i2cDev.Tx([]byte{value | pinE}, nil)
	time.Sleep(1 * time.Microsecond)

	// Send data with E low (falling edge latches data)
	i2cDev.Tx([]byte{value}, nil)
	time.Sleep(50 * time.Microsecond)
}

// writeCommand sends a command byte to the LCD
func writeCommand(cmd byte) {
	// Send high nibble (RS=0 for command).
	write4bits(cmd & 0xF0)
	// Send low nibble.
	write4bits((cmd << 4) & 0xF0)
}

// writeData sends a data byte to the LCD.
func writeData(data byte) {
	// Send high nibble (RS=1 for data).
	write4bits((data & 0xF0) | pinRS)
	// Send low nibble.
	write4bits(((data << 4) & 0xF0) | pinRS)
}

func Display(msg string) {
	Clear()
	fmt.Fprint(os.Stdout, msg+"\n")
	BacklightOn()
	Message(msg)
}

func Clear() {
	BacklightOff()
	if i2cDev != nil {
		writeCommand(cmdClear)
		time.Sleep(2 * time.Millisecond)
	}
}

func BacklightOn() {
	backlight = pinBacklight
	if i2cDev != nil {
		i2cDev.Tx([]byte{backlight}, nil)
	}
}

func BacklightOff() {
	backlight = 0x00
	if i2cDev != nil {
		i2cDev.Tx([]byte{backlight}, nil)
	}
}

// Message displays text on the LCD (supports 2 lines separated by \n)
func Message(text string) {
	if i2cDev == nil {
		return
	}

	lines := strings.Split(text, "\n")

	// DDRAM addresses for 16x2 LCD
	// Line 0: 0x00-0x0F (0x80-0x8F with command bit)
	// Line 1: 0x40-0x4F (0xC0-0xCF with command bit)
	rowOffsets := []byte{0x00, 0x40}

	for row := 0; row < Rows && row < len(lines); row++ {
		line := lines[row]

		// Truncate to column width
		if len(line) > Cols {
			line = line[:Cols]
		}

		// Set cursor to beginning of line
		writeCommand(cmdSetDDRAMAddr | rowOffsets[row])

		// Write each character
		for i := 0; i < len(line); i++ {
			writeData(line[i])
		}
	}
}
