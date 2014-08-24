package main

import (
	"fmt"
	"image/color"
	"time"

	twodee "../libs/twodee"
)

type HudLayer struct {
	text        *twodee.TextRenderer
	regularFont *twodee.FontFace
	planetFont  *twodee.FontFace
	globalText  *twodee.TextCache
	tempText    map[int]*twodee.TextCache
	popText     map[int]*twodee.TextCache
	bounds      twodee.Rectangle
	App         *Application
	game        *GameLayer
}

func NewHudLayer(app *Application, game *GameLayer) (layer *HudLayer, err error) {
	var (
		regularFont *twodee.FontFace
		planetFont  *twodee.FontFace
		background  = color.Transparent
		font        = "assets/fonts/Exo-SemiBold.ttf"
	)
	if regularFont, err = twodee.NewFontFace(font, 24, color.RGBA{255, 255, 255, 255}, background); err != nil {
		return
	}
	if planetFont, err = twodee.NewFontFace(font, 18, color.RGBA{255, 255, 255, 255}, background); err != nil {
		return
	}
	layer = &HudLayer{
		regularFont: regularFont,
		planetFont:  planetFont,
		tempText:    map[int]*twodee.TextCache{},
		popText:     map[int]*twodee.TextCache{},
		globalText:  twodee.NewTextCache(regularFont),
		App:         app,
		bounds:      twodee.Rect(0, 0, 1024, 768),
		game:        game,
	}
	err = layer.Reset()
	return
}

func (l *HudLayer) Delete() {
	if l.text != nil {
		l.text.Delete()
	}
	for _, v := range l.tempText {
		v.Delete()
	}
	for _, v := range l.popText {
		v.Delete()
	}
	l.globalText.Delete()
}

func (l *HudLayer) Render() {
	var (
		textCache     *twodee.TextCache
		planetPos     twodee.Point
		screenPos     twodee.Point
		adjust        twodee.Point
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
	text = fmt.Sprintf("POPULATION: %d     RECORD: %d", aggPopulation, maxPopulation)
	l.globalText.SetText(text)
	if l.globalText.Texture != nil {
		y = maxY - float32(l.globalText.Texture.Height)
		l.text.Draw(l.globalText.Texture, 5, y)
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
		planetPos = planet.Pos()
		if textCache, ok = l.popText[p]; !ok {
			textCache = twodee.NewTextCache(l.planetFont)
			l.popText[p] = textCache
		}
		textCache.SetText(fmt.Sprintf("%d PEOPLE", planet.GetPopulation()))
		if textCache.Texture != nil {
			adjust = twodee.Pt(planet.Radius+0.1, -planet.Radius-0.1)
			screenPos = l.game.WorldToScreenCoords(planetPos.Add(adjust))
			l.text.Draw(textCache.Texture, screenPos.X, screenPos.Y-float32(textCache.Texture.Height))
		}
		//Display Individual Planet Temperatures
		if textCache, ok = l.tempText[p]; !ok {
			textCache = twodee.NewTextCache(l.regularFont)
			l.tempText[p] = textCache
		}
		textCache.SetText(fmt.Sprintf("%d°F", planet.GetTemperature()))
		if textCache.Texture != nil {
			adjust = twodee.Pt(planet.Radius+0.1, planet.Radius+0.1)
			screenPos = l.game.WorldToScreenCoords(planetPos.Add(adjust))
			l.text.Draw(textCache.Texture, screenPos.X, screenPos.Y)
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
	l.Delete()
	if l.text, err = twodee.NewTextRenderer(l.bounds); err != nil {
		return
	}
	return
}
