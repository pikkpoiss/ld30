package main

import (
	"../libs/twodee"
	"time"
)

type OverlayLayer struct {
	Game               *GameLayer
	Events             *twodee.GameEventHandler
	TileRenderer       *twodee.TileRenderer
	Bounds             twodee.Rectangle
	Showing            bool
	Frame              int
	observerShowSplash int
}

func NewOverlayLayer(app *Application, game *GameLayer) (layer *OverlayLayer, err error) {
	layer = &OverlayLayer{
		Game:   game,
		Events: app.GameEventHandler,
		Bounds: twodee.Rect(-48, -36, 48, 36),
	}
	tilem := twodee.TileMetadata{
		Path:       "/../assets/sun.psd",
		PxPerUnit:  320,
		TileWidth:  320,
		TileHeight: 320,
		FramesWide: 1,
		FramesHigh: 1,
	}
	if layer.TileRenderer, err = twodee.NewTileRenderer(layer.Bounds, app.WinBounds, tilem); err != nil {
		return
	}
	layer.observerShowSplash = layer.Events.AddObserver(ShowEndScreen, layer.OnShowEndScreen)
	return
}

func (l *OverlayLayer) OnShowEndScreen(e twodee.GETyper) {
}

func (l *OverlayLayer) Delete() {
	if l.TileRenderer != nil {
		l.TileRenderer.Delete()
	}
	l.Events.RemoveObserver(ShowEndScreen, l.observerShowSplash)
}

func (l *OverlayLayer) Show(frame int) {
	l.Showing = true
	l.Frame = frame
}

func (l *OverlayLayer) Render() {
	if !l.Showing {
		return
	}
	l.TileRenderer.Bind()
	l.TileRenderer.Draw(l.Frame, 0.5, 0.5, 0, false, false)
	l.TileRenderer.Unbind()
}

func (l *OverlayLayer) HandleEvent(evt twodee.Event) bool {
	if !l.Showing {
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
