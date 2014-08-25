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
	tileM            twodee.TileMetadata
	text             *twodee.TextRenderer
	regFont          *twodee.FontFace
	offset           twodee.Point
	popFont          *twodee.FontFace
	maxPopCache      *twodee.TextCache
	cheevosCache     map[int]*twodee.TextCache
	visible          bool
	frame            int
	gameOverObserver int
}

func NewOverlayLayer(app *Application, game *GameLayer) (layer *OverlayLayer, err error) {
	var (
		regFont *twodee.FontFace
		popFont *twodee.FontFace
		bg      = color.Transparent
		exoFont = "assets/fonts/Exo-SemiBold.ttf"
		tileM   twodee.TileMetadata
	)
	if regFont, err = twodee.NewFontFace(exoFont, 24, regColor, bg); err != nil {
		return
	}
	if popFont, err = twodee.NewFontFace(exoFont, 24, hiColor, bg); err != nil {
		return
	}
	tileM = twodee.TileMetadata{
		Path:       "assets/final.png",
		PxPerUnit:  1,
		TileWidth:  1024,
		TileHeight: 768,
		FramesWide: 1,
		FramesHigh: 1,
	}
	layer = &OverlayLayer{
		app:          app,
		game:         game,
		events:       app.GameEventHandler,
		bounds:       app.WinBounds,
		regFont:      regFont,
		popFont:      popFont,
		offset:       twodee.Pt(200, 260),
		cheevosCache: map[int]*twodee.TextCache{},
		maxPopCache:  twodee.NewTextCache(popFont),
		visible:      false,
		frame:        overlayEndFrame,
		tileM:        tileM,
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
	l.maxPopCache.Delete()
	for _, v := range l.cheevosCache {
		v.Clear()
	}
	l.events.RemoveObserver(GameOver, l.gameOverObserver)
}

func (l *OverlayLayer) Render() {
	if !l.visible {
		return
	}
	var (
		y         = l.bounds.Max.Y - l.offset.Y
		x         = l.offset.X
		texture   *twodee.Texture
		textCache *twodee.TextCache
		ok        bool
	)
	l.tileRenderer.Bind()
	l.tileRenderer.Draw(l.frame, 512, 384, 0, false, false)
	l.tileRenderer.Unbind()

	l.text.Bind()
	l.maxPopCache.SetText(fmt.Sprintf("Maximum Population: %d", l.game.Sim.GetMaxPopulation()))
	if l.maxPopCache.Texture != nil {
		y = y - float32(l.maxPopCache.Texture.Height)
		l.text.Draw(l.maxPopCache.Texture, x, y)
	}
	// Put a little padding between max pop and cheevos.
	y -= 15
	for i, item := range l.game.Cheevos.Passed {
		if textCache, ok = l.cheevosCache[i]; !ok {
			textCache = twodee.NewTextCache(l.regFont)
			l.cheevosCache[i] = textCache
		}
		textCache.SetText(item)
		texture = textCache.Texture
		if texture != nil {
			y = y - float32(texture.Height)
			l.text.Draw(texture, x, y)
		}
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
