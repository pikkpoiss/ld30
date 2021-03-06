package main

import (
	"math"
	"time"

	twodee "../libs/twodee"
)

const (
	magicVelocityScalingFactor = 1e-3
	// Starting time is 5minutes.
	startDur = time.Duration(5) * time.Minute
)

type GameLayer struct {
	BatchRenderer         *twodee.BatchRenderer
	TileRenderer          *twodee.TileRenderer
	GlowRenderer          *GlowRenderer
	Bounds                twodee.Rectangle
	App                   *Application
	Sim                   *Simulation
	Starmap               *twodee.Batch
	Cheevos               *Cheevos
	MouseX                float32
	MouseY                float32
	DropPlanetListener    int
	ReleasePlanetListener int
	openMenuListener      int
	closeMenuListener     int
	gameOverListener      int
	DurLeft               time.Duration
	phantomPlanet         *PlanetaryBody
	count                 int64
	paused                bool
}

func NewGameLayer(app *Application) (layer *GameLayer, err error) {
	var bounds = twodee.Rect(-48, -36, 48, 36)
	layer = &GameLayer{
		App:           app,
		Bounds:        bounds,
		Sim:           NewSimulation(bounds, app.GameEventHandler),
		DurLeft:       startDur,
		phantomPlanet: nil,
		count:         0,
		paused:        false,
	}
	if layer.BatchRenderer, err = twodee.NewBatchRenderer(layer.Bounds, app.WinBounds); err != nil {
		return
	}
	tilem := twodee.TileMetadata{
		Path:          "assets/tiles.png",
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
	layer.Cheevos = NewCheevos(layer.App.GameEventHandler, layer.Sim)
	layer.DropPlanetListener = layer.App.GameEventHandler.AddObserver(DropPlanet, layer.OnDropPlanet)
	layer.ReleasePlanetListener = layer.App.GameEventHandler.AddObserver(ReleasePlanet, layer.OnReleasePlanet)
	layer.openMenuListener = layer.App.GameEventHandler.AddObserver(MenuOpen, layer.OnMenuToggle)
	layer.closeMenuListener = layer.App.GameEventHandler.AddObserver(MenuClose, layer.OnMenuToggle)
	layer.gameOverListener = layer.App.GameEventHandler.AddObserver(GameOver, layer.OnGameOver)
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
	if l.Cheevos != nil {
		l.Cheevos.Delete()
	}
	l.App.GameEventHandler.RemoveObserver(DropPlanet, l.DropPlanetListener)
	l.App.GameEventHandler.RemoveObserver(ReleasePlanet, l.ReleasePlanetListener)
	l.App.GameEventHandler.RemoveObserver(MenuOpen, l.openMenuListener)
	l.App.GameEventHandler.RemoveObserver(MenuClose, l.closeMenuListener)
	l.App.GameEventHandler.RemoveObserver(GameOver, l.gameOverListener)
}

func (l *GameLayer) Render() {
	var (
		pos     twodee.Point
		radians float64
	)
	l.count = (l.count + 2) % 100000000
	radians = 0.0174532925 * float64(l.count)
	var glow = math.Sin(radians) * 0.1
	l.GlowRenderer.Bind()
	l.GlowRenderer.SetStrength(float32(0.3 + glow))

	/*
		// Enable planets to block out aura.
		l.TileRenderer.Bind()
		l.GlowRenderer.DisableOutput()
		for _, p := range l.Sim.Planets {
			pos = p.Pos()
			l.TileRenderer.DrawScaled(p.Frame(), pos.X, pos.Y, 0, p.Scale, false, false)
		}
		l.GlowRenderer.EnableOutput()
		l.TileRenderer.Unbind()
	*/

	l.BatchRenderer.Bind()
	l.BatchRenderer.Draw(l.Starmap, l.Bounds.Min.X, l.Bounds.Min.Y, 0)
	l.BatchRenderer.Unbind()

	l.TileRenderer.Bind()
	// Draw an aura around planets.
	for _, p := range l.Sim.Planets {
		pos = p.Pos()
		frame := 33
		if p.HasState(Dying) {
			frame = p.Frame()
		}
		l.TileRenderer.DrawScaled(frame, pos.X, pos.Y, 0, p.Scale, false, false)
	}
	pos = l.Sim.Sun.Pos()
	l.TileRenderer.Draw(l.Sim.Sun.Frame(), pos.X, pos.Y, float32(radians), false, false)

	l.GlowRenderer.Unbind()

	l.TileRenderer.Unbind()

	l.BatchRenderer.Bind()
	l.BatchRenderer.Draw(l.Starmap, l.Bounds.Min.X, l.Bounds.Min.Y, 0)
	l.BatchRenderer.Unbind()

	l.TileRenderer.Bind()
	l.TileRenderer.Draw(l.Sim.Sun.Frame(), pos.X, pos.Y, float32(radians), false, false)
	for _, p := range l.Sim.Planets {
		pos = p.Pos()
		l.TileRenderer.DrawScaled(p.Frame(), pos.X, pos.Y, p.Rotation, p.Scale, false, false)
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
	if l.paused {
		return
	}
	l.Sim.Update(elapsed)
	l.Cheevos.Update(elapsed)
	l.DurLeft -= elapsed
	if l.DurLeft <= 0 {
		l.DurLeft = time.Duration(0)
		l.App.GameEventHandler.Enqueue(twodee.NewBasicGameEvent(GameOver))
	}
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
			l.App.GameEventHandler.Enqueue(NewReleasePlanetEvent(l.MouseX, l.MouseY))
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
			relVector := twodee.Pt(event.X-p.X, event.Y-p.Y)
			// Since the vector's magnitude is still too big, we
			// need to scale it down by some magic factor.
			relVector = relVector.Scale(magicVelocityScalingFactor)
			l.phantomPlanet.Velocity = relVector
			l.phantomPlanet.RemState(Phantom)
			l.Sim.AddPlanet(l.phantomPlanet)
			l.phantomPlanet = nil
		}
	}
}

func (l *GameLayer) OnMenuToggle(evt twodee.GETyper) {
	l.paused = !l.paused
}

func (l *GameLayer) OnGameOver(evt twodee.GETyper) {
	l.paused = true
}

func (l *GameLayer) WorldToScreenCoords(pt twodee.Point) twodee.Point {
	x, y := l.TileRenderer.WorldToScreenCoords(pt.X, pt.Y)
	return twodee.Pt(x, y)
}
