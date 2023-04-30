package main

import (
	"time"

	"github.com/faiface/pixel"
)

type Daemon struct {
	Matrix pixel.Matrix

	Health uint8

	AnimTicker *time.Ticker
	Sprites    []*pixel.Sprite
	AnimStep   uint8 // next drawing
}

func NewDaemon(s []*pixel.Sprite, m pixel.Matrix) *Daemon {
	at := time.NewTicker(1000 / 12 * time.Millisecond)
	return &Daemon{
		Matrix:     m,
		Health:     20,
		AnimTicker: at,
		Sprites:    s,
		AnimStep:   0,
	}
}

func (d *Daemon) Draw(t pixel.Target) {
	d.Sprites[d.AnimStep].Draw(t, d.Matrix)
}

func (d *Daemon) Update(tick int) error {
	d.Matrix = d.Matrix.Moved(pixel.V(1.4, 0))
	select {
	case <-d.AnimTicker.C:
		d.AnimStep++
		d.AnimStep %= uint8(len(d.Sprites))

	default:
	}
	return nil
}

func (d *Daemon) Hit(dmg uint8) {
	d.Health -= dmg
}

func (d *Daemon) IsDead() bool {
	return d.Health <= 0
}
