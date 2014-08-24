package main

import (
	"time"

	twodee "../libs/twodee"
)

type GameLayer struct {
	BatchRenderer         *twodee.BatchRenderer
	TileRenderer          *twodee.TileRenderer
	GlowRenderer          *GlowRenderer
	Bounds                twodee.Rectangle
	App                   *Application
	Sim                   *Simulation
	Starmap               *twodee.Batch
	MouseX                float32
	MouseY                float32
	DropPlanetListener    int
	ReleasePlanetListener int
	phantomPlanet         *PlanetaryBody
}

func NewGameLayer(app *Application) (layer *GameLayer, err error) {
	var bounds = twodee.Rect(-32, -24, 32, 24)
	layer = &GameLayer{
		App:           app,
		Bounds:        bounds,
		Sim:           NewSimulation(bounds, app.GameEventHandler),
		phantomPlanet: nil,
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
	if layer.GlowRenderer, err = NewGlowRenderer(192, 128, 6, 0.3, 1.0); err != nil {
		return
	}
	if layer.Starmap, err = LoadMap("assets/starmap.tmx"); err != nil {
		return
	}
	layer.DropPlanetListener = layer.App.GameEventHandler.AddObserver(DropPlanet, layer.OnDropPlanet)
	layer.ReleasePlanetListener = layer.App.GameEventHandler.AddObserver(ReleasePlanet, layer.OnReleasePlanet)
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
	l.App.GameEventHandler.RemoveObserver(ReleasePlanet, l.ReleasePlanetListener)
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
		l.TileRenderer.DrawScaled(p.Frame(), pos.X, pos.Y, 0, p.Scale, false, false)
	}
	// TODO: Render l.phantomPlanet.
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
		l.TileRenderer.DrawScaled(p.Frame(), pos.X, pos.Y, 0, p.Scale, false, false)
	}
	if l.phantomPlanet != nil {
		p := l.phantomPlanet
		pos = p.Pos()
		l.TileRenderer.DrawScaled(p.Frame(), pos.X, pos.Y, 0, p.Scale, false, false)
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
		switch event.Type {
		case twodee.Press:
			l.App.GameEventHandler.Enqueue(NewDropPlanetEvent(l.MouseX, l.MouseY))
		case twodee.Release:
			var mag float32 = 0.001
			l.App.GameEventHandler.Enqueue(NewReleasePlanetEvent(l.MouseX, l.MouseY, mag))
		default:
			break
		}
	case *twodee.MouseMoveEvent:
		l.MouseX, l.MouseY = l.TileRenderer.ScreenToWorldCoords(event.X, event.Y)
	}
	return true
}

func (l *GameLayer) OnDropPlanet(evt twodee.GETyper) {
	switch event := evt.(type) {
	case *DropPlanetEvent:
		l.phantomPlanet = NewPlanet(event.X, event.Y)
		l.phantomPlanet.SetState(Phantom)
	}
}

func (l *GameLayer) OnReleasePlanet(evt twodee.GETyper) {
	switch event := evt.(type) {
	case *ReleasePlanetEvent:
		// do something.
		if l.phantomPlanet != nil {
			p := l.phantomPlanet.Pos()
			relVector := twodee.Pt(event.P.X-p.X, event.P.Y-p.Y)
			// TODO: I'm not really sure we actually need to scale
			// the relative vector by a magnitude, since we get
			// that for free by virtue of the relativeness of the
			// vector to some point p0.
			// Still, we need to reduce it by some amount so the
			// planet doesn't go careening off screen immediately.
			relVector = relVector.Scale(event.Mag)
			l.phantomPlanet.Velocity = relVector
			l.phantomPlanet.RemState(Phantom)
			l.Sim.AddPlanet(l.phantomPlanet)
			l.phantomPlanet = nil
		}
	}
}

func (l *GameLayer) WorldToScreenCoords(pt twodee.Point) twodee.Point {
	x, y := l.TileRenderer.WorldToScreenCoords(pt.X, pt.Y)
	return twodee.Pt(x, y)
}
