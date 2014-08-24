package main

import (
	"fmt"
	"image/color"
	"strconv"
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

func NewHudLayer(app *Application, game *GameLayer) (layer *HudLayer, err error) {
	var (
		regularFont *twodee.FontFace
		background  = color.Transparent
		font        = "assets/fonts/Roboto-Black.ttf"
	)
	if regularFont, err = twodee.NewFontFace(font, 24, color.RGBA{255, 255, 255, 255}, background); err != nil {
		return
	}
	layer = &HudLayer{
		regularFont: regularFont,
		cache:       map[int]*twodee.TextCache{},
		App:         app,
		bounds:      twodee.Rect(0, 0, 1024, 768),
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
		textCache     *twodee.TextCache
		texture       *twodee.Texture
		ok            bool
		text          string
		x, y          float32
		maxX          = l.bounds.Max.X
		maxY          = l.bounds.Max.Y
		aggPopulation = l.game.Sim.GetPopulation()
		maxPopulation = l.game.Sim.GetMaxPopulation()
	)
	l.text.Bind()

	// Display Aggregate Population Count
	if textCache, ok = l.cache[1]; !ok {
		textCache = twodee.NewTextCache(l.regularFont)
		l.cache[1] = textCache
	}
	text = fmt.Sprintf("POPULATION: %d     RECORD: %d", aggPopulation, maxPopulation)
	textCache.SetText(text)
	texture = textCache.Texture
	if texture != nil {
		y = maxY - float32(texture.Height)
		l.text.Draw(texture, 5, y)
	}

	// Display time remaining.
	s := int64(l.game.DurLeft.Seconds())
	m := s / 60
	s = s % 60
	text = fmt.Sprintf("%d:%2d", m, s)
	textCache.SetText(text)
	texture = textCache.Texture
	if texture != nil {
		y = maxY - float32(texture.Height)
		// Some fudged padding to make sure there's room for the clock.
		x = maxX - 60.0
		l.text.Draw(texture, x, y)
	}

	//Display Individual Planet Population Counts
	for p, planet := range l.game.Sim.Planets {
		if textCache, ok = l.cache[p]; !ok {
			textCache = twodee.NewTextCache(l.regularFont)
			l.cache[p] = textCache
		}
		pos := planet.Pos()
		textCache.SetText(strconv.Itoa(planet.GetPopulation()))
		texture = textCache.Texture
		if texture != nil {
			adjust := twodee.Pt(planet.Radius+0.1, -planet.Radius-1.75)
			pt := l.game.WorldToScreenCoords(pos.Add(adjust))
			l.text.Draw(texture, pt.X, pt.Y)
		}
		//Display Individual Planet Temperatures
		textCache.SetText(strconv.Itoa(int(planet.GetTemperature())) + "°F")
		texture = textCache.Texture
		if texture != nil {
			adjust := twodee.Pt(planet.Radius+0.1, planet.Radius)
			pt := l.game.WorldToScreenCoords(pos.Add(adjust))
			l.text.Draw(texture, pt.X, pt.Y)
		}
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
