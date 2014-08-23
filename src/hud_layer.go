package main

import (
	"image/color"
	"time"

	twodee "../libs/twodee"
)

type HudLayer struct {
	text        *twodee.TextRenderer
	regularFont *twodee.FontFace
	cache       map[int]*twodee.TextCache
	bounds      twodee.Rectangle
	App         *Application
	game        *GameLayer
}

const (
	HudHeight = 1200
	HudWidth  = 1200
)

func NewHudLayer(app *Application, game *GameLayer) (layer *HudLayer, err error) {
	var (
		regularFont *twodee.FontFace
		background  = color.Transparent
		font        = "assets/fonts/Roboto-Black.ttf"
	)
	if regularFont, err = twodee.NewFontFace(font, 32, color.RGBA{255, 255, 255, 255}, background); err != nil {
		return
	}
	layer = &HudLayer{
		regularFont: regularFont,
		cache:       map[int]*twodee.TextCache{},
		App:         app,
		bounds:      twodee.Rect(0, 0, HudWidth, HudHeight),
		game:        game,
	}
	err = layer.Reset()
	return
}

func (l *HudLayer) Delete() {
	l.text.Delete()
	for _, v := range l.cache {
		v.Delete()
	}
}

func (l *HudLayer) Render() {
	var (
		textCache *twodee.TextCache
		texture   *twodee.Texture
		ok        bool
		y         = l.bounds.Max.Y
	)
	l.text.Bind()
	if textCache, ok = l.cache[1]; !ok {
		textCache = twodee.NewTextCache(l.regularFont)
		l.cache[1] = textCache
	}
	textCache.SetText("POPULATION: 0")
	texture = textCache.Texture
	if texture != nil {
		y = y - float32(texture.Height)
		l.text.Draw(texture, 0, y)
	}
	l.text.Unbind()
}

func (l *HudLayer) HandleEvent(evt twodee.Event) bool {
	return true
}

func (l *HudLayer) Update(elapsed time.Duration) {
}

func (l *HudLayer) Reset() (err error) {
	if l.text != nil {
		l.text.Delete()
	}
	if l.text, err = twodee.NewTextRenderer(l.bounds); err != nil {
		return
	}
	for _, v := range l.cache {
		v.Delete()
	}
	return
}
