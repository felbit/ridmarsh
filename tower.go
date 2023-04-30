package main

import (
	"time"

	"github.com/faiface/pixel"
)

type Missile struct {
	Matrix    pixel.Matrix
	Direction pixel.Vec
	Speed     float64
}

type Tower struct {
	Matrix pixel.Matrix
	Sprite *pixel.Sprite

	Target     *Daemon
	ShotTicker *time.Ticker
	Range      float64
}

func NewTower(s *pixel.Sprite, m pixel.Matrix) *Tower {
	return &Tower{
		Matrix:     m,
		Sprite:     s,
		ShotTicker: time.NewTicker(time.Second),
	}
}

func (t *Tower) Draw(target pixel.Target) {
	t.Sprite.Draw(target, t.Matrix)
}

func (t *Tower) Shoot() bool {
	select {
	case <-t.ShotTicker.C:
		return true
	default:
		return false
	}
}
