package main

import (
	"../libs/twodee"
	"time"
)

type HudLayer struct {
	Bounds twodee.Rectangle
	App    *Application
	game   *GameLayer
}

const (
	HudHeight = 100
	HudWidth  = 500
)

func NewHudLayer(app *Application, game *GameLayer) (layer *HudLayer, err error) {
	layer = &HudLayer{
		App:    app,
		Bounds: twodee.Rect(0, 0, HudWidth, HudHeight),
		game:   game,
	}
	return
}

func (l *HudLayer) Delete() {
}

func (l *HudLayer) Render() {
}

func (l *HudLayer) HandleEvent(evt twodee.Event) bool {
	return true
}

func (l *HudLayer) Update(elapsed time.Duration) {
}

func (l *HudLayer) Reset() (err error) {
	return
}
