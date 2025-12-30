package lcd

import (
	"fmt"
	"os"
	"strings"
	"time"

	"periph.io/x/conn/v3/gpio"
	"periph.io/x/conn/v3/gpio/gpioreg"
	"periph.io/x/host/v3"
)

// Pin number = GPIO number.
var (
	pinRs        = 7
	pinEn        = 8
	pinD4        = 25
	pinD5        = 24
	pinD6        = 23
	pinD7        = 17
	pinBacklight = 11

	lcd  *display
	Cols = 16
)

// HD44780 command constants
const (
	cmdClear        = 0x01
	cmdHome         = 0x02
	cmdEntryMode    = 0x04
	cmdDisplayCtrl  = 0x08
	cmdCursorShift  = 0x10
	cmdFunctionSet  = 0x20
	cmdSetCGRAMAddr = 0x40
	cmdSetDDRAMAddr = 0x80

	// Entry mode flags
	entryLeft          = 0x02
	entryShiftDecrement = 0x00

	// Display control flags
	displayOn  = 0x04
	displayOff = 0x00
	cursorOn   = 0x02
	cursorOff  = 0x00
	blinkOn    = 0x01
	blinkOff   = 0x00

	// Function set flags
	func4BitMode = 0x00
	func8BitMode = 0x10
	func1Line    = 0x00
	func2Line    = 0x08
	func5x8Dots  = 0x00
	func5x10Dots = 0x04

	// Timing constants
	writeDelay = 37 * time.Microsecond
	pulseDelay = 1 * time.Microsecond
	clearDelay = 2 * time.Millisecond
)

// Row address offsets for 16-column displays
var rowOffsets = []byte{0x00, 0x40}

type hd44780 struct {
	rs, en, d4, d5, d6, d7, backlight gpio.PinOut
	rows                               int
	cols                               int
}

type display struct {
	hd   *hd44780
	cols int
	rows int
}

func Stop() {
	if lcd != nil {
		lcd.Clear()
	}
	// periph.io handles cleanup automatically
}

func Start() {
	// Initialize periph.io host drivers
	if _, err := host.Init(); err != nil {
		panic(err)
	}

	// Get GPIO pins
	rsPin := gpioreg.ByName(fmt.Sprintf("GPIO%d", pinRs))
	enPin := gpioreg.ByName(fmt.Sprintf("GPIO%d", pinEn))
	d4Pin := gpioreg.ByName(fmt.Sprintf("GPIO%d", pinD4))
	d5Pin := gpioreg.ByName(fmt.Sprintf("GPIO%d", pinD5))
	d6Pin := gpioreg.ByName(fmt.Sprintf("GPIO%d", pinD6))
	d7Pin := gpioreg.ByName(fmt.Sprintf("GPIO%d", pinD7))
	blPin := gpioreg.ByName(fmt.Sprintf("GPIO%d", pinBacklight))

	if rsPin == nil || enPin == nil || d4Pin == nil ||
		d5Pin == nil || d6Pin == nil || d7Pin == nil || blPin == nil {
		panic("failed to get GPIO pins")
	}

	// Initialize HD44780 controller
	hd := newHD44780(rsPin, enPin, d4Pin, d5Pin, d6Pin, d7Pin, blPin, 2, Cols)

	// Create display wrapper
	lcd = &display{
		hd:   hd,
		cols: Cols,
		rows: 2,
	}

	lcd.Clear()
}

func Display(msg string) {
	lcd.Clear()
	fmt.Fprint(os.Stdout, msg+"\n")
	lcd.BacklightOn()
	lcd.Message(msg)
}

func Clear() {
	lcd.BacklightOff()
	lcd.Clear()
}

// newHD44780 creates a new HD44780 controller in 4-bit GPIO mode
func newHD44780(rs, en, d4, d5, d6, d7, backlight gpio.PinOut, rows, cols int) *hd44780 {
	hd := &hd44780{
		rs:        rs,
		en:        en,
		d4:        d4,
		d5:        d5,
		d6:        d6,
		d7:        d7,
		backlight: backlight,
		rows:      rows,
		cols:      cols,
	}

	// Initialize all pins to low
	hd.rs.Out(gpio.Low)
	hd.en.Out(gpio.Low)
	hd.d4.Out(gpio.Low)
	hd.d5.Out(gpio.Low)
	hd.d6.Out(gpio.Low)
	hd.d7.Out(gpio.Low)
	hd.backlight.Out(gpio.Low)

	// HD44780 initialization sequence for 4-bit mode
	time.Sleep(50 * time.Millisecond) // Wait for LCD to power up

	// Send 0x03 three times (initialization)
	hd.write4bits(0x03)
	time.Sleep(5 * time.Millisecond)
	hd.write4bits(0x03)
	time.Sleep(150 * time.Microsecond)
	hd.write4bits(0x03)
	time.Sleep(150 * time.Microsecond)

	// Switch to 4-bit mode
	hd.write4bits(0x02)
	time.Sleep(150 * time.Microsecond)

	// Function set: 4-bit mode, 2 lines, 5x8 dots
	hd.writeCommand(cmdFunctionSet | func4BitMode | func2Line | func5x8Dots)

	// Display control: display on, cursor off, blink off
	hd.writeCommand(cmdDisplayCtrl | displayOn | cursorOff | blinkOff)

	// Entry mode: left to right, no shift
	hd.writeCommand(cmdEntryMode | entryLeft | entryShiftDecrement)

	// Clear display
	hd.Clear()

	return hd
}

// write4bits writes a 4-bit value to the data pins
func (hd *hd44780) write4bits(value byte) {
	// Set data pins
	if value&0x01 != 0 {
		hd.d4.Out(gpio.High)
	} else {
		hd.d4.Out(gpio.Low)
	}
	if value&0x02 != 0 {
		hd.d5.Out(gpio.High)
	} else {
		hd.d5.Out(gpio.Low)
	}
	if value&0x04 != 0 {
		hd.d6.Out(gpio.High)
	} else {
		hd.d6.Out(gpio.Low)
	}
	if value&0x08 != 0 {
		hd.d7.Out(gpio.High)
	} else {
		hd.d7.Out(gpio.Low)
	}

	// Pulse enable pin to latch data
	time.Sleep(pulseDelay)
	hd.en.Out(gpio.High)
	time.Sleep(pulseDelay)
	hd.en.Out(gpio.Low)
	time.Sleep(pulseDelay)
}

// writeByte writes an 8-bit value as two 4-bit values
func (hd *hd44780) writeByte(value byte, mode gpio.Level) {
	hd.rs.Out(mode)

	// Write high nibble
	hd.write4bits(value >> 4)

	// Write low nibble
	hd.write4bits(value & 0x0F)

	time.Sleep(writeDelay)
}

// writeCommand writes a command to the LCD
func (hd *hd44780) writeCommand(cmd byte) {
	hd.writeByte(cmd, gpio.Low)
}

// writeChar writes a character to the LCD
func (hd *hd44780) writeChar(char byte) {
	hd.writeByte(char, gpio.High)
}

// Clear clears the display
func (hd *hd44780) Clear() {
	hd.writeCommand(cmdClear)
	time.Sleep(clearDelay)
}

// Home returns cursor to home position
func (hd *hd44780) Home() {
	hd.writeCommand(cmdHome)
	time.Sleep(clearDelay)
}

// SetCursor sets the cursor position
func (hd *hd44780) SetCursor(col, row int) {
	if row >= len(rowOffsets) {
		row = len(rowOffsets) - 1
	}
	hd.writeCommand(cmdSetDDRAMAddr | (byte(col) + rowOffsets[row]))
}

// BacklightOn turns the backlight on
func (hd *hd44780) BacklightOn() {
	hd.backlight.Out(gpio.High)
}

// BacklightOff turns the backlight off
func (hd *hd44780) BacklightOff() {
	hd.backlight.Out(gpio.Low)
}

// Display methods

// Clear clears the display and turns off backlight
func (d *display) Clear() {
	d.hd.Clear()
}

// BacklightOn turns the backlight on
func (d *display) BacklightOn() {
	d.hd.BacklightOn()
}

// BacklightOff turns the backlight off
func (d *display) BacklightOff() {
	d.hd.BacklightOff()
}

// Message writes a message to the display, handling newlines
func (d *display) Message(text string) {
	lines := strings.Split(text, "\n")

	for row := 0; row < d.rows && row < len(lines); row++ {
		d.hd.SetCursor(0, row)
		line := lines[row]

		// Truncate or pad to column width
		if len(line) > d.cols {
			line = line[:d.cols]
		}

		// Write each character
		for _, char := range line {
			d.hd.writeChar(byte(char))
		}
	}
}
