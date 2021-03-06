package lcd

import (
	"fmt"
	"os"

	"github.com/kidoman/embd"
	"github.com/kidoman/embd/controller/hd44780"
	_ "github.com/kidoman/embd/host/rpi"
	"github.com/kidoman/embd/interface/display/characterdisplay"
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

	lcd  *characterdisplay.Display
	Cols = 16
)

func Stop() {
	lcd.Clear()
	embd.CloseGPIO()
}

func Start() {
	if err := embd.InitGPIO(); err != nil {
		panic(err)
	}

	hd, err := hd44780.NewGPIO(
		pinRs,
		pinEn,
		pinD4,
		pinD5,
		pinD6,
		pinD7,
		pinBacklight,
		false,
		hd44780.RowAddress16Col,
		hd44780.TwoLine,
	)
	if err != nil {
		panic(err)
	}

	lcd = characterdisplay.New(hd, Cols, 2)
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
