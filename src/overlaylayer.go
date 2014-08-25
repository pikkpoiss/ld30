package main

import (
	"fmt"
	"image/color"
	"time"

	"../libs/twodee"
)

const (
	overlayEndFrame = 0
)

type OverlayLayer struct {
	game             *GameLayer
	app              *Application
	events           *twodee.GameEventHandler
	tileRenderer     *twodee.TileRenderer
	bounds           twodee.Rectangle
	text             *twodee.TextRenderer
	regFont          *twodee.FontFace
	tileM            twodee.TileMetadata
	cheevosText      *twodee.TextCache
	visible          bool
	frame            int
	gameOverObserver int
}

func NewOverlayLayer(app *Application, game *GameLayer) (layer *OverlayLayer, err error) {
	var (
		regFont *twodee.FontFace
		bg      = color.Transparent
		exoFont = "assets/fonts/Exo-SemiBold.ttf"
		tileM   twodee.TileMetadata
	)
	if regFont, err = twodee.NewFontFace(exoFont, 24, regColor, bg); err != nil {
		return
	}
	tileM = twodee.TileMetadata{
		// TODO: Get an overlay png full of cool tiles.
		Path:       "assets/sun.png",
		PxPerUnit:  320,
		TileWidth:  320,
		TileHeight: 320,
		FramesWide: 1,
		FramesHigh: 1,
	}
	layer = &OverlayLayer{
		app:         app,
		game:        game,
		events:      app.GameEventHandler,
		bounds:      twodee.Rect(-48, -36, 48, 36),
		regFont:     regFont,
		cheevosText: twodee.NewTextCache(regFont),
		visible:     false,
		tileM:       tileM,
	}
	layer.Reset()
	return
}

func (l *OverlayLayer) OnGameOver(e twodee.GETyper) {
	l.visible = true
}

func (l *OverlayLayer) Delete() {
	l.visible = false
	if l.tileRenderer != nil {
		l.tileRenderer.Delete()
	}
	l.cheevosText.Delete()
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
	var (
		y    float32
		maxY = l.bounds.Max.Y
	)
	l.tileRenderer.Bind()
	l.tileRenderer.Draw(l.frame, 0.5, 0.5, 0, false, false)
	l.tileRenderer.Unbind()

	l.text.Bind()
	text := fmt.Sprintf("hihi")
	l.cheevosText.SetText(text)
	if l.cheevosText.Texture != nil {
		y = maxY - float32(l.cheevosText.Texture.Height)
		l.text.Draw(l.cheevosText.Texture, 5, y)
	}
	l.text.Unbind()
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
			return l.NewGame()
		}
	}
	// Handle all events.
	return false
}

func (l *OverlayLayer) NewGame() bool {
	if err := l.Reset(); err != nil {
		// TODO: Ugh, see if we can reset GameLayer?
		return false
	}
	return false
}

func (l *OverlayLayer) Update(elapsed time.Duration) {
}

func (l *OverlayLayer) Reset() (err error) {
	l.Delete()
	if l.tileRenderer, err = twodee.NewTileRenderer(l.bounds, l.app.WinBounds, l.tileM); err != nil {
		return
	}
	l.gameOverObserver = l.events.AddObserver(GameOver, l.OnGameOver)
	if l.text, err = twodee.NewTextRenderer(l.bounds); err != nil {
		return
	}
	return
}
