package main

import (
	"fmt"
	"image/color"
	"time"

	twodee "../libs/twodee"
)

type HudLayer struct {
	text            *twodee.TextRenderer
	regularFont     *twodee.FontFace
	planetFont      *twodee.FontFace
	messageFont     *twodee.FontFace
	messageText     *twodee.TextCache
	messageCoords   twodee.Point
	globalText      *twodee.TextCache
	timeText        *twodee.TextCache
	tempText        map[int]*twodee.TextCache
	popText         map[int]*twodee.TextCache
	bounds          twodee.Rectangle
	App             *Application
	game            *GameLayer
	messageListener int
}

func NewHudLayer(app *Application, game *GameLayer) (layer *HudLayer, err error) {
	var (
		regularFont *twodee.FontFace
		planetFont  *twodee.FontFace
		messageFont *twodee.FontFace
		background  = color.Transparent
		exoFont     = "assets/fonts/Exo-SemiBold.ttf"
		abelFont    = "assets/fonts/Abel-Regular.ttf"
	)
	if regularFont, err = twodee.NewFontFace(exoFont, 24, regColor, background); err != nil {
		return
	}
	if planetFont, err = twodee.NewFontFace(exoFont, 18, regColor, background); err != nil {
		return
	}
	if messageFont, err = twodee.NewFontFace(abelFont, 40, regColor, background); err != nil {
		return
	}
	layer = &HudLayer{
		regularFont: regularFont,
		planetFont:  planetFont,
		messageFont: messageFont,
		tempText:    map[int]*twodee.TextCache{},
		popText:     map[int]*twodee.TextCache{},
		globalText:  twodee.NewTextCache(regularFont),
		timeText:    twodee.NewTextCache(regularFont),
		messageText: twodee.NewTextCache(messageFont),
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
	l.timeText.Delete()
	l.messageText.Delete()
	l.App.GameEventHandler.RemoveObserver(DisplayMessage, l.messageListener)
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
	if m > 0 {
		text = fmt.Sprintf("%d:%02d", m, s)
	} else {
		text = fmt.Sprintf("%02d", s)
	}
	l.timeText.SetText(text)
	if l.timeText.Texture != nil {
		y = maxY - float32(l.timeText.Texture.Height)
		// Some fudged padding to make sure there's room for the clock.
		x = maxX - 60.0
		l.text.Draw(l.timeText.Texture, x, y)
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
		textCache.SetText(fmt.Sprintf("%v %d°F", planet.Name, planet.GetTemperature()))
		if textCache.Texture != nil {
			adjust = twodee.Pt(planet.Radius+0.1, planet.Radius+0.1)
			screenPos = l.game.WorldToScreenCoords(planetPos.Add(adjust))
			l.text.Draw(textCache.Texture, screenPos.X, screenPos.Y)
		}
	}
	if l.messageText.Texture != nil {
		l.text.Draw(l.messageText.Texture, l.messageCoords.X, l.messageCoords.Y)
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
	l.messageListener = l.App.GameEventHandler.AddObserver(DisplayMessage, l.OnDisplayMessage)
	return
}

func (l *HudLayer) OnDisplayMessage(evt twodee.GETyper) {
	switch event := evt.(type) {
	case *DisplayMessageEvent:
		l.messageText.SetText(event.Message)
		if l.messageText.Texture != nil {
			if event.Positioned {
				l.messageCoords = l.game.WorldToScreenCoords(event.Coords)
			} else {
				l.messageCoords = twodee.Point{
					(l.bounds.Max.X - float32(l.messageText.Texture.OriginalWidth)) / 2.0,
					(l.bounds.Max.Y - float32(l.messageText.Texture.OriginalHeight)) / 4.0,
				}
			}
		}
	}
}
