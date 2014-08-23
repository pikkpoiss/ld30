package main

import (
	twodee "../libs/twodee"
	"github.com/go-gl/gl"
	"runtime"
	"time"
)

func init() {
	// See https://code.google.com/p/go/issues/detail?id=3527
	runtime.LockOSThread()
}

type Application struct {
	layers                *twodee.Layers
	Context               *twodee.Context
	WinBounds             twodee.Rectangle
	GameEventHandler      *twodee.GameEventHandler
	gameClosingObserverId int
	InitiateCloseGame     bool
}

func NewApplication() (app *Application, err error) {
	var (
		layers            *twodee.Layers
		context           *twodee.Context
		gameLayer         *GameLayer
		winbounds         = twodee.Rect(0, 0, 1024, 768)
		gameEventHandler  = twodee.NewGameEventHandler(NumGameEventTypes)
		initiateCloseGame = false
	)
	if context, err = twodee.NewContext(); err != nil {
		return
	}
	context.SetFullscreen(false)
	context.SetCursor(true)
	if err = context.CreateWindow(int(winbounds.Max.X), int(winbounds.Max.Y), "LD30"); err != nil {
		return
	}
	layers = twodee.NewLayers()
	app = &Application{
		layers:            layers,
		Context:           context,
		WinBounds:         winbounds,
		GameEventHandler:  gameEventHandler,
		InitiateCloseGame: initiateCloseGame,
	}
	if gameLayer, err = NewGameLayer(app); err != nil {
		return
	}
	layers.Push(gameLayer)
	app.gameClosingObserverId = app.GameEventHandler.AddObserver(GameIsClosing, app.CloseGame)
	return
}

func (a *Application) Draw() {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	a.layers.Render()
}

func (a *Application) Update(elapsed time.Duration) {
	a.layers.Update(elapsed)
}

func (a *Application) Delete() {
	a.GameEventHandler.RemoveObserver(GameIsClosing, a.gameClosingObserverId)
	a.layers.Delete()
	a.Context.Delete()
}

func (a *Application) ProcessEvents() {
	var (
		evt  twodee.Event
		loop = true
	)
	for loop {
		select {
		case evt = <-a.Context.Events.Events:
			a.layers.HandleEvent(evt)
		default:
			// No more events
			loop = false
		}
	}
}

func (a *Application) CloseGame(e twodee.GETyper) {
	a.InitiateCloseGame = true
}

func main() {
	var (
		app *Application
		err error
	)

	if app, err = NewApplication(); err != nil {
		panic(err)
	}
	defer app.Delete()

	var (
		last_render  = time.Now()
		current_time = time.Now()
		updated_to   = current_time
		step         = twodee.Step60Hz
		render_max   = step
	)
	for !app.Context.Window.ShouldClose() && !app.InitiateCloseGame {
		for !updated_to.After(current_time) {
			app.Update(step)
			updated_to = updated_to.Add(step)
		}
		if diff := current_time.Sub(last_render); diff < render_max {
			time.Sleep(render_max - diff)
		}
		app.Draw()
		app.Context.Window.SwapBuffers()
		last_render = current_time
		app.Context.Events.Poll()
		app.GameEventHandler.Poll()
		app.ProcessEvents()
		current_time = time.Now()
	}
}
