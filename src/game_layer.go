package main

import (
	twodee "../libs/twodee"
	"time"
)

type GameLayer struct {
	BatchRenderer *twodee.BatchRenderer
	TileRenderer  *twodee.TileRenderer
	Bounds        twodee.Rectangle
	App           *Application
}

func NewGameLayer(app *Application) (layer *GameLayer, err error) {
	layer = &GameLayer{
		App:    app,
		Bounds: twodee.Rect(0, 0, 20, 20),
	}
	if layer.BatchRenderer, err = twodee.NewBatchRenderer(layer.Bounds, app.WinBounds); err != nil {
		return
	}
	tilem := twodee.TileMetadata{
		Path:       "assets/tiles.fw.png",
		PxPerUnit:  int(PxPerUnit),
		TileWidth:  32,
		TileHeight: 32,
	}
	if layer.TileRenderer, err = twodee.NewTileRenderer(layer.Bounds, app.WinBounds, tilem); err != nil {
		return
	}
	return
}

func (l *GameLayer) Delete() {
	if l.TileRenderer != nil {
		l.TileRenderer.Delete()
	}
	if l.BatchRenderer != nil {
		l.BatchRenderer.Delete()
	}
}

func (l *GameLayer) Render() {
	l.BatchRenderer.Bind()
	l.BatchRenderer.Bind()
	l.TileRenderer.Bind()
	l.TileRenderer.Unbind()
	return
}

func (l *GameLayer) Update(elapsed time.Duration) {
}

func (l *GameLayer) Reset() (err error) {
	return
}

func (l *GameLayer) HandleEvent(evt twodee.Event) bool {
	return true
}
