package main

import (
	"image/color"
	"strconv"
	"time"

	"tinygo.org/x/drivers/lis3dh"
	"tinygo.org/x/drivers/scd4x"
	"tinygo.org/x/drivers/st7789"
	"tinygo.org/x/drivers/ws2812"
	"tinygo.org/x/tinyfont"
	"tinygo.org/x/tinyfont/freesans"

	"machine"
)

var display st7789.Device
var leds ws2812.Device
var accel lis3dh.Device
var bzrPin machine.Pin
var btnA, btnB, btnUp, btnLeft, btnDown, btnRight machine.Pin

const (
	BLACK = iota
	WHITE
	RED
	SNAKE
	TEXT
	ORANGE
	PURPLE
)

var colors = []color.RGBA{
	color.RGBA{0, 0, 0, 255},
	color.RGBA{255, 255, 255, 255},
	color.RGBA{250, 0, 0, 255},
	color.RGBA{0, 200, 0, 255},
	color.RGBA{160, 160, 160, 255},
	color.RGBA{255, 153, 51, 255},
	color.RGBA{153, 51, 255, 255},
}

var (
	co2sensor *scd4x.Device
)

func main() {
	time.Sleep(3 * time.Second)

	machine.SPI0.Configure(machine.SPIConfig{
		Frequency: 8000000,
		Mode:      0,
	})

	machine.I2C0.Configure(machine.I2CConfig{
		SCL: machine.I2C0_SCL_PIN,
		SDA: machine.I2C0_SDA_PIN,
	})
	accel = lis3dh.New(machine.I2C0)
	accel.Address = lis3dh.Address0
	accel.Configure()

	display = st7789.New(machine.SPI0,
		machine.TFT_RST,       // TFT_RESET
		machine.TFT_WRX,       // TFT_DC
		machine.TFT_CS,        // TFT_CS
		machine.TFT_BACKLIGHT) // TFT_LITE

	display.Configure(st7789.Config{
		Rotation: st7789.ROTATION_270,
		Height:   320,
	})

	neo := machine.NEOPIXELS
	neo.Configure(machine.PinConfig{Mode: machine.PinOutput})
	leds = ws2812.New(neo)

	bzrPin = machine.SPEAKER
	bzrPin.Configure(machine.PinConfig{Mode: machine.PinOutput})

	speaker := machine.SPEAKER_ENABLE
	speaker.Configure(machine.PinConfig{Mode: machine.PinOutput})
	speaker.High()

	co2sensor = scd4x.New(machine.I2C0)
	co2sensor.Configure()

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

	black := color.RGBA{0, 0, 0, 255}
	display.FillScreen(black)

	x, y, z := accel.ReadRawAcceleration()
	println("XYZ", x, y, z)
	if err := co2sensor.StartPeriodicMeasurement(); err != nil {
		println(err)
	}

	var co2 int32
	var temp int32
	var hum int32
	var err error
	var c color.RGBA
	oldc := color.RGBA{R: 0xff, G: 0xff, B: 0xff}
	for {

		for i := 0; i < 5; i++ {
			co2, err = co2sensor.ReadCO2()
			temp, _ = co2sensor.ReadTemperature()
			println("TEMP", temp)
			hum, _ = co2sensor.ReadHumidity()
			println("HUM", hum)
			if err != nil {
				println(err)
			}
			println(co2)
			if co2 != 0 {
				break
			} else {
				time.Sleep(200 * time.Millisecond)
			}
		}
		switch {
		case co2 < 800:
			c = color.RGBA{R: 0x00, G: 0xff, B: 0x00}
		case co2 < 1500:
			c = color.RGBA{R: 0xff, G: 0xff, B: 0x00}
		default:
			c = color.RGBA{R: 0xff, G: 0x00, B: 0x00}
		}
		if c != oldc {
			display.FillScreen(c)
			oldc = c
		}

		tinyfont.WriteLine(&display, &freesans.Bold12pt7b, 20, 90, "Temperature: "+strconv.FormatFloat(float64(temp/1000), 'f', 2, 64)+" C", black)
		tinyfont.WriteLine(&display, &freesans.Bold12pt7b, 20, 120, "Humidity: "+strconv.Itoa(int(hum))+" %", black)
		tinyfont.WriteLine(&display, &freesans.Bold12pt7b, 20, 180, "CO2 level: "+strconv.Itoa(int(co2))+" ppm", black)

		time.Sleep(1 * time.Second)

		tinyfont.WriteLine(&display, &freesans.Bold12pt7b, 20, 90, "Temperature: "+strconv.FormatFloat(float64(temp/1000), 'f', 2, 64)+" C", c)
		tinyfont.WriteLine(&display, &freesans.Bold12pt7b, 20, 120, "Humidity: "+strconv.Itoa(int(hum))+" %", c)
		tinyfont.WriteLine(&display, &freesans.Bold12pt7b, 20, 180, "CO2 level: "+strconv.Itoa(int(co2))+" ppm", c)
	}

}
