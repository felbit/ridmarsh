package main

import (
	"fmt"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font/basicfont"
)

type UI struct {
	atlas     *text.Atlas
	gameStats *text.Text

	health, economy, fiendCount int
}

func NewUI() *UI {
	orig := pixel.V(
		-float64(gridW*gridS/2-10),
		float64(gridH*gridS/2-30),
	)
	atlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)
	label := text.New(orig, atlas)
	label.Color = colornames.Antiquewhite

	loudMessage := text.New(pixel.ZV, atlas)
	loudMessage.Color = colornames.Crimson
	return &UI{atlas: atlas, gameStats: label}
}

func (ui *UI) Update(health, economy, fiends int) {
	ui.health = health
	ui.economy = economy
	ui.fiendCount = fiends
}

func (ui *UI) Draw(t pixel.Target) {
	ui.gameStats.Clear()
	fmt.Fprintf(ui.gameStats, "H: %d \n", ui.health)
	fmt.Fprintf(ui.gameStats, "$: %d \n", ui.economy)
	fmt.Fprintf(ui.gameStats, "F: %d \n", ui.fiendCount)
	ui.gameStats.Draw(t, pixel.IM.Scaled(ui.gameStats.Orig, 2))
}

func (ui *UI) GameOver(t pixel.Target) {
	gameOver := text.New(pixel.ZV, ui.atlas)
	gameOver.Color = colornames.Ghostwhite

	fmt.Fprintln(gameOver, "GAME OVER!")
	gameOver.Draw(t, pixel.IM.Scaled(gameOver.Orig, 5).Moved(pixel.V(screenW/2-gameOver.Bounds().Max.X*6, -128)))

	restart := text.New(pixel.ZV, ui.atlas)
	restart.Color = colornames.Crimson

	fmt.Fprintln(restart, "press ENTER to restart")
	restart.Draw(t, pixel.IM.Scaled(restart.Orig, 2).Moved(pixel.V(screenW/2-restart.Bounds().Max.X*2, -160)))
}
