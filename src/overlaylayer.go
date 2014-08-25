package main

import (
	"time"

	"../libs/twodee"
)

type OverlayLayer struct {
	game             *GameLayer
	events           *twodee.GameEventHandler
	tileRenderer     *twodee.TileRenderer
	bounds           twodee.Rectangle
	visible          bool
	frame            int
	gameOverObserver int
}

func NewOverlayLayer(app *Application, game *GameLayer) (layer *OverlayLayer, err error) {
	layer = &OverlayLayer{
		game:    game,
		events:  app.GameEventHandler,
		bounds:  twodee.Rect(-48, -36, 48, 36),
		visible: false,
	}
	tilem := twodee.TileMetadata{
		Path:       "/../assets/sun.psd",
		PxPerUnit:  320,
		TileWidth:  320,
		TileHeight: 320,
		FramesWide: 1,
		FramesHigh: 1,
	}
	if layer.tileRenderer, err = twodee.NewTileRenderer(layer.bounds, app.WinBounds, tilem); err != nil {
		return
	}
	layer.gameOverObserver = layer.events.AddObserver(GameOver, layer.OnGameOver)
	return
}

func (l *OverlayLayer) OnGameOver(e twodee.GETyper) {
	l.visible = true
}

func (l *OverlayLayer) Delete() {
	if l.tileRenderer != nil {
		l.tileRenderer.Delete()
	}
	l.events.RemoveObserver(GameOver, l.gameOverObserver)
}

func (l *OverlayLayer) Show(frame int) {
	l.visible = true
	l.frame = frame
}

func (l *OverlayLayer) Render() {
	if !l.visible {
		return
	}
	l.tileRenderer.Bind()
	l.tileRenderer.Draw(l.frame, 0.5, 0.5, 0, false, false)
	l.tileRenderer.Unbind()
}

func (l *OverlayLayer) HandleEvent(evt twodee.Event) bool {
	if !l.visible {
		return true
	}
	switch event := evt.(type) {
	case *twodee.KeyEvent:
		if event.Type != twodee.Press {
			break
		}
		switch event.Code {
		case twodee.KeyEnter:
			return l.Advance()
		}
	}
	return true
}

func (l *OverlayLayer) Advance() bool {
	return false
}

func (l *OverlayLayer) Update(elapsed time.Duration) {
}

func (l *OverlayLayer) Reset() (err error) {
	return
}
