package main

import (
	"log"
	"math/rand"
	"time"

	"github.com/faiface/pixel"
)

type Fiend struct {
	Matrix pixel.Matrix

	Health float64

	AnimTicker *time.Ticker
	Sprites    []*pixel.Sprite
	AnimStep   uint8 // next drawing
}

func NewFiend() *Fiend {
	// random location
	y := float64(rand.Intn(gridH)-gridH/2) * float64(gridS)
	log.Println("Y:", y)
	sprites := make([]*pixel.Sprite, 0, 6)
	for _, r := range AnimationFiendRun {
		sprites = append(sprites, pixel.NewSprite(FiendRaw, r))
	}
	at := time.NewTicker(1000 / 12 * time.Millisecond)
	return &Fiend{
		Matrix:     pixel.IM.Moved(pixel.V(-float64(gridW*gridS/2), y)),
		Health:     20,
		AnimTicker: at,
		Sprites:    sprites,
		AnimStep:   0,
	}
}

func (f *Fiend) Draw(t pixel.Target) {
	f.Sprites[f.AnimStep].Draw(t, f.Matrix)
}

func (f *Fiend) Update(tick uint64) error {
	f.Matrix = f.Matrix.Moved(pixel.V(1.4, 0))
	select {
	case <-f.AnimTicker.C:
		f.AnimStep++
		f.AnimStep %= uint8(len(f.Sprites))

	default:
	}
	return nil
}

func (f *Fiend) Position() pixel.Vec {
	return pixel.Vec{f.Matrix[4], f.Matrix[5]}
}

func (f *Fiend) Hit(dmg float64) {
	f.Health -= dmg
}

func (f *Fiend) IsDead() bool {
	return f.Health <= 0
}
