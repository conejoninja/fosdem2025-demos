package main

import (
	"image/color"
	"machine"
	"time"

	"tinygo.org/x/drivers/st7789"
)

var display st7789.Device
var btnA, btnB, btnUp, btnLeft, btnDown, btnRight machine.Pin

func main() {
	// Setup the SPI connection of the GopherBadge
	machine.SPI0.Configure(machine.SPIConfig{
		Frequency: 8000000,
		Mode:      0,
	})

	// Create a new display device
	display = st7789.New(machine.SPI0,
		machine.TFT_RST,       // TFT_RESET
		machine.TFT_WRX,       // TFT_DC
		machine.TFT_CS,        // TFT_CS
		machine.TFT_BACKLIGHT) // TFT_LITE

	display.Configure(st7789.Config{
		Rotation: st7789.ROTATION_270,
		Height:   320,
	})

	// Setup the buttons
	btnA = machine.BUTTON_A
	btnB = machine.BUTTON_B
	btnUp = machine.BUTTON_UP
	btnLeft = machine.BUTTON_LEFT
	btnDown = machine.BUTTON_DOWN
	btnRight = machine.BUTTON_RIGHT
	btnA.Configure(machine.PinConfig{Mode: machine.PinInput})
	btnB.Configure(machine.PinConfig{Mode: machine.PinInput})
	btnUp.Configure(machine.PinConfig{Mode: machine.PinInput})
	btnLeft.Configure(machine.PinConfig{Mode: machine.PinInput})
	btnDown.Configure(machine.PinConfig{Mode: machine.PinInput})
	btnRight.Configure(machine.PinConfig{Mode: machine.PinInput})

	green := color.RGBA{0, 255, 0, 255}
	black := color.RGBA{0, 0, 0, 255}

	// fill the whole screen with black
	display.FillScreen(black)

	w, h := display.Size()
	x := (w - 10) / 2
	y := (h - 10) / 2
	for {

		// "clear" our previous square position
		display.FillRectangle(x, y, 10, 10, black)

		// Check if any direction button is pressed
		if !btnLeft.Get() {
			x -= 10
		}
		if !btnUp.Get() {
			y -= 10
		}
		if !btnDown.Get() {
			y += 10
		}
		if !btnRight.Get() {
			x += 10
		}

		// draw our little square
		display.FillRectangle(x, y, 10, 10, green)

		time.Sleep(100 * time.Millisecond)
	}
}
