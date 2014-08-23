package main

import (
	twodee "../libs/twodee"
)

const (
	GameIsClosing twodee.GameEventType = iota
	DropPlanet
	sentinel
)

const (
	NumGameEventTypes = int(sentinel)
)

type DropPlanetEvent struct {
	twodee.BasicGameEvent
	X float32
	Y float32
}

func NewDropPlanetEvent(x, y float32) (e *DropPlanetEvent) {
	e = &DropPlanetEvent{
		*twodee.NewBasicGameEvent(DropPlanet),
		x,
		y,
	}
	return
}
