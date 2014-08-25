package main

import (
	twodee "../libs/twodee"
)

const (
	GameIsClosing twodee.GameEventType = iota
	PlayBackgroundMusic
	DropPlanet
	ReleasePlanet
	ShowSplash
	PauseMusic
	PlayGameOverEffect
	DisplayMessage
	MenuOpen
	MenuClose
	MenuClick
	MenuSel
	sentinel
)

const (
	NumGameEventTypes = int(sentinel)
)

type DisplayMessageEvent struct {
	twodee.BasicGameEvent
	X       float32
	Y       float32
	Message string
}

type DropPlanetEvent struct {
	twodee.BasicGameEvent
	X float32
	Y float32
}

type ReleasePlanetEvent DropPlanetEvent

func NewDropPlanetEvent(x, y float32) (e *DropPlanetEvent) {
	e = &DropPlanetEvent{
		*twodee.NewBasicGameEvent(DropPlanet),
		x,
		y,
	}
	return
}

func NewReleasePlanetEvent(x, y float32) (e *ReleasePlanetEvent) {
	e = &ReleasePlanetEvent{
		*twodee.NewBasicGameEvent(ReleasePlanet),
		x,
		y,
	}
	return
}

func NewDisplayMessageEvent(x, y float32, message string) (e *DisplayMessageEvent) {
	e = &DisplayMessageEvent{
		*twodee.NewBasicGameEvent(DisplayMessage),
		x,
		y,
		message,
	}
	return
}
