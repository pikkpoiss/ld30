package main

import (
	twodee "../libs/twodee"
)

const (
	GameIsClosing twodee.GameEventType = iota
	PlayBackgroundMusic
	DropPlanet
	ReleasePlanet
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

type ReleasePlanetEvent struct {
	twodee.BasicGameEvent
	P   twodee.Point
	Mag float32
}

func NewDropPlanetEvent(x, y float32) (e *DropPlanetEvent) {
	e = &DropPlanetEvent{
		*twodee.NewBasicGameEvent(DropPlanet),
		x,
		y,
	}
	return
}

func NewReleasePlanetEvent(x, y, m float32) (e *ReleasePlanetEvent) {
	e = &ReleasePlanetEvent{
		*twodee.NewBasicGameEvent(ReleasePlanet),
		twodee.Pt(x, y),
		m,
	}
	return
}
