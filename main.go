package main

import (
	"image"
	_ "image/png"
	"log"
	"math"
	"os"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

func Main() {
	cfg := pixelgl.WindowConfig{
		Title:     "Ridmarsh - A Ravensong Saga",
		Bounds:    pixel.R(0, 0, 900, 450),
		VSync:     true,
		Resizable: true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}
	win.SetSmooth(true)

	// canvas := pixelgl.NewCanvas(pixel.R(-320, -240, 320, 240))

	towerRaw, _ := loadPicture("assets/tower_left.png")
	daemonRaw, _ := loadPicture("assets/cacodaemon.png")

	var DaemonAnimRun []pixel.Rect
	for x := daemonRaw.Bounds().Min.X; x < 6*64; x += 64 {
		DaemonAnimRun = append(DaemonAnimRun, pixel.R(x, 3*64, x+64, 4*64))
	}

	var (
		camPos       = pixel.ZV
		camZoom      = 1.0
		camZoomSpeed = 1.2
		towers       []*Tower
		daemons      []*Daemon
	)

	tick := 0
	last := time.Now()
	for !win.Closed() {
		delta := time.Since(last).Seconds()
		last = time.Now()
		_ = delta

		cam := pixel.IM.Scaled(camPos, camZoom).Moved(win.Bounds().Center().Sub(camPos))
		win.SetMatrix(cam)

		camZoom *= math.Pow(camZoomSpeed, win.MouseScroll().Y)

		if win.JustPressed(pixelgl.MouseButtonLeft) {
			mouse := cam.Unproject(win.MousePosition())
			tower := NewTower(
				pixel.NewSprite(towerRaw, towerRaw.Bounds()),
				pixel.IM.Moved(mouse),
			)
			towers = append(towers, tower)
		}

		if win.JustPressed(pixelgl.MouseButtonRight) {
			mouse := cam.Unproject(win.MousePosition())
			sprites := make([]*pixel.Sprite, 0, 6)
			for _, r := range DaemonAnimRun {
				sprites = append(sprites, pixel.NewSprite(daemonRaw, r))
			}
			daemon := NewDaemon(
				sprites,
				pixel.IM.Moved(mouse),
			)
			daemons = append(daemons, daemon)
		}

		win.Clear(colornames.Greenyellow)

		for _, tower := range towers {
			tpos := pixel.V(tower.Matrix[4], tower.Matrix[5])
			for _, d := range daemons {
				dpos := pixel.V(d.Matrix[4], d.Matrix[5])
				dist := distance(tpos, dpos)
				if dist < 150 && tower.Shoot() {
					log.Println("BOOM!")
					d.Hit(10)
				}
			}
			tower.Draw(win)
		}

		for i, daemon := range daemons {
			if daemon.IsDead() {
				ds := daemons[:i]
				if len(daemons) > i+1 {
					ds = append(ds, daemons[i+1:]...)
				}
				daemons = ds
				continue
			}
			daemon.Update(tick)
			daemon.Draw(win)
		}
		win.Update()
		tick++
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
