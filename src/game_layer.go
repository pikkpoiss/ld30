package main

import (
	twodee "../libs/twodee"
	"time"
)

type GameLayer struct {
	BatchRenderer      *twodee.BatchRenderer
	TileRenderer       *twodee.TileRenderer
	Bounds             twodee.Rectangle
	App                *Application
	Sim                *Simulation
	MouseX             float32
	MouseY             float32
	DropPlanetListener int
}

func NewGameLayer(app *Application) (layer *GameLayer, err error) {
	layer = &GameLayer{
		App:    app,
		Bounds: twodee.Rect(-28, -21, 28, 21),
		Sim:    NewSimulation(),
	}
	if layer.BatchRenderer, err = twodee.NewBatchRenderer(layer.Bounds, app.WinBounds); err != nil {
		return
	}
	tilem := twodee.TileMetadata{
		Path:       "assets/tiles.fw.png",
		PxPerUnit:  int(PxPerUnit),
		TileWidth:  128,
		TileHeight: 128,
	}
	if layer.TileRenderer, err = twodee.NewTileRenderer(layer.Bounds, app.WinBounds, tilem); err != nil {
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
	l.App.GameEventHandler.RemoveObserver(DropPlanet, l.DropPlanetListener)
}

func (l *GameLayer) Render() {
	l.BatchRenderer.Bind()
	l.BatchRenderer.Bind()
	l.TileRenderer.Bind()
	pos := l.Sim.Sun.Pos()
	l.TileRenderer.Draw(l.Sim.Sun.Frame(), pos.X, pos.Y, 0, false, false)
	for _, p := range l.Sim.Planets {
		pos = p.Pos()
		l.TileRenderer.Draw(p.Frame(), pos.X, pos.Y, 0, false, false)
	}
	l.TileRenderer.Unbind()
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
