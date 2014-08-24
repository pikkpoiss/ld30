package main

import (
	twodee "../libs/twodee"
	"time"
)

type GameLayer struct {
	BatchRenderer      *twodee.BatchRenderer
	TileRenderer       *twodee.TileRenderer
	GlowRenderer       *GlowRenderer
	Bounds             twodee.Rectangle
	App                *Application
	Sim                *Simulation
	Starmap            *twodee.Batch
	MouseX             float32
	MouseY             float32
	DropPlanetListener int
}

func NewGameLayer(app *Application) (layer *GameLayer, err error) {
	var bounds = twodee.Rect(-28, -21, 28, 21)
	layer = &GameLayer{
		App:    app,
		Bounds: bounds,
		Sim:    NewSimulation(bounds, app.GameEventHandler),
	}
	if layer.BatchRenderer, err = twodee.NewBatchRenderer(layer.Bounds, app.WinBounds); err != nil {
		return
	}
	tilem := twodee.TileMetadata{
		Path:          "assets/tiles.fw.png",
		PxPerUnit:     int(PxPerUnit),
		TileWidth:     128,
		TileHeight:    128,
		Interpolation: twodee.NearestInterpolation,
	}
	if layer.TileRenderer, err = twodee.NewTileRenderer(layer.Bounds, app.WinBounds, tilem); err != nil {
		return
	}
	if layer.GlowRenderer, err = NewGlowRenderer(192, 128, 8, 0.3, 1.0); err != nil {
		return
	}
	if layer.Starmap, err = LoadMap("assets/starmap.tmx"); err != nil {
		return
	}
	layer.DropPlanetListener = layer.App.GameEventHandler.AddObserver(DropPlanet, layer.OnDropPlanet)
	return
}

func (l *GameLayer) Delete() {
	if l.TileRenderer != nil {
		l.TileRenderer.Delete()
	}
	if l.BatchRenderer != nil {
		l.BatchRenderer.Delete()
	}
	if l.GlowRenderer != nil {
		l.GlowRenderer.Delete()
	}
	if l.Starmap != nil {
		l.Starmap.Delete()
	}
	l.App.GameEventHandler.RemoveObserver(DropPlanet, l.DropPlanetListener)
}

func (l *GameLayer) Render() {
	var (
		pos twodee.Point
	)
	l.GlowRenderer.Bind()

	l.TileRenderer.Bind()
	l.GlowRenderer.DisableOutput()
	for _, p := range l.Sim.Planets {
		pos = p.Pos()
		l.TileRenderer.Draw(p.Frame(), pos.X, pos.Y, 0, false, false)
	}
	l.GlowRenderer.EnableOutput()
	l.TileRenderer.Unbind()

	l.BatchRenderer.Bind()
	l.BatchRenderer.Draw(l.Starmap, l.Bounds.Min.X, l.Bounds.Min.Y, 0)
	l.BatchRenderer.Unbind()

	l.TileRenderer.Bind()
	pos = l.Sim.Sun.Pos()
	l.TileRenderer.Draw(l.Sim.Sun.Frame(), pos.X, pos.Y, 0, false, false)

	l.GlowRenderer.Unbind()

	l.TileRenderer.Unbind()

	l.BatchRenderer.Bind()
	l.BatchRenderer.Draw(l.Starmap, l.Bounds.Min.X, l.Bounds.Min.Y, 0)
	l.BatchRenderer.Unbind()

	l.TileRenderer.Bind()
	l.TileRenderer.Draw(l.Sim.Sun.Frame(), pos.X, pos.Y, 0, false, false)
	for _, p := range l.Sim.Planets {
		pos = p.Pos()
		l.TileRenderer.Draw(p.Frame(), pos.X, pos.Y, 0, false, false)
	}
	l.TileRenderer.Unbind()

	l.GlowRenderer.Draw()
	return
}

func (l *GameLayer) Update(elapsed time.Duration) {
	l.Sim.Update(elapsed)
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
	case *twodee.MouseButtonEvent:
		if event.Type != twodee.Press {
			break
		}
		l.App.GameEventHandler.Enqueue(NewDropPlanetEvent(l.MouseX, l.MouseY))
		l.App.GameEventHandler.Enqueue(twodee.NewBasicGameEvent(PlayPlanetDropEffect))
	case *twodee.MouseMoveEvent:
		l.MouseX, l.MouseY = l.TileRenderer.ScreenToWorldCoords(event.X, event.Y)
	}
	return true
}

func (l *GameLayer) OnDropPlanet(evt twodee.GETyper) {
	switch event := evt.(type) {
	case *DropPlanetEvent:
		l.Sim.AddPlanet(event.X, event.Y)
	}
}

func (l *GameLayer) WorldToScreenCoords(pt twodee.Point) twodee.Point {
	x, y := l.TileRenderer.WorldToScreenCoords(pt.X, pt.Y)
	return twodee.Pt(x, y)
}
