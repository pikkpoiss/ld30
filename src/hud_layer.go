package main

import (
	twodee "../libs/twodee"
	"fmt"
	"image/color"
	"time"
)

type HudLayer struct {
	text        *twodee.TextRenderer
	regularFont *twodee.FontFace
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
		background  = color.Transparent
		font        = "assets/fonts/Roboto-Black.ttf"
	)
	if regularFont, err = twodee.NewFontFace(font, 24, color.RGBA{255, 255, 255, 255}, background); err != nil {
		return
	}
	layer = &HudLayer{
		regularFont: regularFont,
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
		y             = l.bounds.Max.Y
		aggPopulation = l.game.Sim.GetPopulation()
		maxPopulation = l.game.Sim.GetMaxPopulation()
	)
	l.text.Bind()

	// Display Aggregate Population Count
	text = fmt.Sprintf("POPULATION: %d     RECORD: %d", aggPopulation, maxPopulation)
	l.globalText.SetText(text)
	if l.globalText.Texture != nil {
		y = y - float32(l.globalText.Texture.Height)
		l.text.Draw(l.globalText.Texture, 0, y)
	}

	//Display Individual Planet Population Counts
	for p, planet := range l.game.Sim.Planets {
		planetPos = planet.Pos()
		if textCache, ok = l.popText[p]; !ok {
			textCache = twodee.NewTextCache(l.regularFont)
			l.popText[p] = textCache
		}
		textCache.SetText(fmt.Sprintf("%d", planet.GetPopulation()))
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
