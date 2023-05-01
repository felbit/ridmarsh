package main

import (
	"image"
	_ "image/png"
	"log"
	"math"
	"os"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

const (
	gridW int = 16
	gridH int = 9
	gridS int = 64

	StartEconomy uint32 = 300
	StartHealth  uint8  = 1
)

var (
	screenW, screenH  float64
	gameOver          bool
	cam               pixel.Matrix
	FiendRaw          pixel.Picture
	TowerRaw          pixel.Picture
	AnimationFiendRun []pixel.Rect
)

// will be called before main()
func init() {
	TowerRaw, _ = loadPicture("assets/tower_left.png")
	FiendRaw, _ = loadPicture("assets/cacodaemon.png")

	for x := FiendRaw.Bounds().Min.X; x < 6*64; x += 64 {
		AnimationFiendRun = append(AnimationFiendRun, pixel.R(x, 3*64, x+64, 4*64))
	}
}

func Main() {
	screenW = float64(gridW * gridS)
	screenH = float64(gridH * gridS)

	cfg := pixelgl.WindowConfig{
		Title:     "Ridmarsh TD",
		Bounds:    pixel.R(0, 0, screenW, screenH),
		VSync:     true,
		Resizable: true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}
	win.SetSmooth(true)

	ui := NewUI()

	canvas := pixelgl.NewCanvas(pixel.R(
		-screenW/2, -screenH/2,
		screenW/2, screenH/2,
	))

	var (
		towers []*Tower
	)

	game := NewGame()
	gameOver = false

	// last := time.Now()
	for !win.Closed() {
		// delta := time.Since(last).Seconds()
		// last = time.Now()

		win.Clear(colornames.Forestgreen)
		if gameOver {
			canvas.Clear(colornames.Crimson)
			ui.GameOver(canvas)
			win.SetMatrix(cam)
			canvas.Draw(win, pixel.IM.Moved(canvas.Bounds().Center()))
			win.Update()

			if win.JustPressed(pixelgl.KeyEnter) {
				game = NewGame()
				gameOver = false
			}

			continue
		}

		canvas.Clear(colornames.Forestgreen)

		if win.JustPressed(pixelgl.MouseButtonLeft) {
			if game.TrySpend(100) {
				mouse := cam.Unproject(win.MousePosition())
				tower := NewTower(
					pixel.NewSprite(TowerRaw, TowerRaw.Bounds()),
					pixel.IM.Moved(mouse),
				)
				towers = append(towers, tower)
			}
		}

		for _, tower := range towers {
			tpos := pixel.V(tower.Matrix[4], tower.Matrix[5])
			for _, fiend := range game.fiends {
				fpos := pixel.V(fiend.Matrix[4], fiend.Matrix[5])
				dist := distance(tpos, fpos)
				if dist < 150 && tower.Shoot() {
					log.Println("BOOM!")
					fiend.Hit(10)
				}
			}
			tower.Draw(canvas)
		}

		game.Update(win)
		game.Draw(canvas)
		if game.health <= 0 {
			gameOver = true
		}

		ui.Update(int(game.health), int(game.economy), len(game.fiends))
		ui.Draw(canvas)

		// stretch canvas to the window bounds
		cam = pixel.IM.Scaled(
			pixel.ZV,
			math.Min(
				win.Bounds().W()/canvas.Bounds().W(),
				win.Bounds().H()/canvas.Bounds().H(),
			)).Moved(win.Bounds().Center())

		win.Clear(colornames.Forestgreen)
		win.SetMatrix(cam)
		canvas.Draw(win, pixel.IM.Moved(canvas.Bounds().Center()))
		win.Update()

		game.tick++
	}
}

func loadPicture(path string) (pixel.Picture, error) {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		panic(err)
	}
	return pixel.PictureDataFromImage(img), nil
}

func distance(v1, v2 pixel.Vec) float64 {
	return math.Sqrt(math.Pow((v2.Y-v1.Y), 2) + math.Pow((v2.X-v1.X), 2))
}

func main() {
	pixelgl.Run(Main)
}
