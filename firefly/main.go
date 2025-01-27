package main

import "github.com/firefly-zero/firefly-go/firefly"

var (
	tx, ty int
	gopher firefly.Image
	myFont firefly.Font
)

func init() {
	firefly.Boot = boot
	firefly.Update = update
	firefly.Render = render
}

func boot() {
	// Initialize phase
	gopher = firefly.LoadFile("gopher", nil).Image()
	myFont = firefly.LoadFile("font", nil).Font()
	tx = 120
	ty = 100
}

func update() {
	// Update game state based on the rules
	pad, _ := firefly.ReadPad(firefly.Combined)
	switch {
	case pad.DPad().Down:
		ty--
	case pad.DPad().Up:
		ty++
	case pad.DPad().Right:
		tx++
	case pad.DPad().Left:
		tx--
	}
}

func render() {
	// Provide feedback to the players
	firefly.ClearScreen(firefly.ColorWhite)
	firefly.DrawImage(gopher, firefly.Point{X: tx, Y: ty})
	firefly.DrawText("FOSDEM 2025", myFont, firefly.Point{X: 10, Y: 150}, firefly.ColorPurple)
}
