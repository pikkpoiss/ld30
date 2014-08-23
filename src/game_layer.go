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
	Sun           *twodee.AnimatingEntity
}

func NewGameLayer(app *Application) (layer *GameLayer, err error) {
	layer = &GameLayer{
		App:    app,
		Bounds: twodee.Rect(-28, -21, 28, 21),
		Sun:    NewSun(),
	}
	if layer.BatchRenderer, err = twodee.NewBatchRenderer(layer.Bounds, app.WinBounds); err != nil {
		return
	}
	tilem := twodee.TileMetadata{
		Path:       "assets/tiles.fw.png",
		PxPerUnit:  int(PxPerUnit),
		TileWidth:  128,
		TileHeight: 128,
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
	frame := l.Sun.Frame()
	l.TileRenderer.Draw(frame, 0, 0, 0, false, false)
	l.TileRenderer.Unbind()
	return
}

func (l *GameLayer) Update(elapsed time.Duration) {
}

func (l *GameLayer) Reset() (err error) {
	return
}

func (l *GameLayer) HandleEvent(evt twodee.Event) bool {
	switch event := evt.(type) {
	case *twodee.KeyEvent:
		if event.Type != twodee.Press {
			break
		}
		switch event.Code {
		case twodee.KeyEscape:
			l.App.GameEventHandler.Enqueue(twodee.NewBasicGameEvent(GameIsClosing))
			return false
		}
	}
	return true
}
