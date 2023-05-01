package main

import (
	"math/rand"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/google/uuid"
)

type Game struct {
	grid    []*Tower
	fiends  map[uuid.UUID]*Fiend
	tick    uint64
	economy uint32
	health  uint8
}

func NewGame() *Game {
	g := &Game{}
	g.economy = StartEconomy
	g.health = StartHealth
	g.grid = make([]*Tower, 0, gridW*gridH)
	g.fiends = make(map[uuid.UUID]*Fiend, 0)
	return g
}

// Try to spend `amount` of economy; returns `false`, if not enough
// economy and `true` if spending was successful
func (g *Game) TrySpend(amount uint32) bool {
	if amount > g.economy {
		return false
	}
	g.economy -= amount
	return true
}

func (g *Game) Gain(amount uint32) {
	g.economy += amount
}

func (g *Game) Update(win *pixelgl.Window) error {
	// randomly spawn monster
	// random Y value for monster (on grid)

	if p := rand.Intn(1000); p < 20 {
		id := uuid.New()
		g.fiends[id] = NewFiend()

	}

	for id, fiend := range g.fiends {
		// is fiend still visible?
		if fiend.Position().X < float64(gridW*gridS)/2 {
			if fiend.IsDead() {
				g.economy += 10 // TODO: make fiend property
				delete(g.fiends, id)
			} else {
				fiend.Update(g.tick)
			}
		} else { // fiend moved out of the right limit
			g.health--
			delete(g.fiends, id)
		}

	}

	return nil
}

func (g *Game) Draw(t pixel.Target) {
	for _, fiend := range g.fiends {
		fiend.Draw(t)
	}
}
