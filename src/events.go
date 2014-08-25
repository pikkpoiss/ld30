package main

import (
	twodee "../libs/twodee"
)

const (
	GameIsClosing twodee.GameEventType = iota
	PlayBackgroundMusic
	DropPlanet
	PlanetFireDeath
	PlanetCollision
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
	Positioned bool
	Coords     twodee.Point
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

func NewPositionedMessageEvent(pt twodee.Point, message string) (e *DisplayMessageEvent) {
	return &DisplayMessageEvent{
		*twodee.NewBasicGameEvent(DisplayMessage),
		true,
		pt,
		message,
	}
}

func NewMessageEvent(message string) (e *DisplayMessageEvent) {
	return &DisplayMessageEvent{
		*twodee.NewBasicGameEvent(DisplayMessage),
		false,
		twodee.Pt(0, 0),
		message,
	}
}
