package main

import (
	twodee "../libs/twodee"
	"time"
)

type Cheevos struct {
	events *twodee.GameEventHandler
	queue  []Cheevo
	sim    *Simulation
	active Cheevo
}

func NewCheevos(events *twodee.GameEventHandler, sim *Simulation) *Cheevos {
	return &Cheevos{
		events: events,
		queue: []Cheevo{
			NewMakeFirstPlanet(),
		},
		sim:    sim,
		active: nil,
	}
}

func (c *Cheevos) Update(elapsed time.Duration) {
	if c.active != nil {
		c.active.Update(elapsed)
	} else {
		c.selectActive()
	}
}

func (c *Cheevos) selectActive() {
	if c.active != nil || len(c.queue) == 0 {
		return
	}
	for i, candidate := range c.queue {
		if candidate.IsAvailable(c.sim) {
			c.active = candidate
			c.queue = append(c.queue[:i], c.queue[i+1:]...)
			c.active.Init(c.events)
			return
		}
	}
}

func (c *Cheevos) Delete() {
}

type Cheevo interface {
	Init(events *twodee.GameEventHandler)
	IsAvailable(sim *Simulation) bool
	GetLabel() string
	Update(elapsed time.Duration)
}

type BaseCheevo struct {
	label     string
	callbacks []*Callback
}

func newBaseCheevo(label string) *BaseCheevo {
	return &BaseCheevo{
		label:     label,
		callbacks: []*Callback{},
	}
}

func (c *BaseCheevo) GetLabel() string {
	return c.label
}

type Callback struct {
	Elapsed time.Duration
	Trigger time.Duration
	Func    func()
	Done    bool
}

func CallAfter(duration time.Duration, f func()) *Callback {
	return &Callback{
		Elapsed: 0,
		Trigger: duration,
		Func:    f,
		Done:    false,
	}
}

func (c *Callback) Update(elapsed time.Duration) {
	c.Elapsed += elapsed
	if c.Elapsed >= c.Trigger {
		c.Done = true
		c.Func()
	}
}

func (c *BaseCheevo) After(d time.Duration, f func()) {
	c.callbacks = append(c.callbacks, CallAfter(d, f))
}

func sendMessage(msg string, events *twodee.GameEventHandler) func() {
	return func() {
		events.Enqueue(NewMessageEvent(msg))
	}
}

func (c *BaseCheevo) SendMessages(messages []string, interval time.Duration, events *twodee.GameEventHandler) {
	counter := interval
	for i := 0; i < len(messages); i++ {
		c.After(counter, sendMessage(messages[i], events))
		counter += interval
	}
	c.After(counter, sendMessage("", events))
}

func (c *BaseCheevo) Update(elapsed time.Duration) {
	for i := len(c.callbacks) - 1; i >= 0; i-- {
		c.callbacks[i].Update(elapsed)
		if c.callbacks[i].Done {
			c.callbacks = append(c.callbacks[:i], c.callbacks[i+1:]...)
		}
	}
}

type MakeFirstPlanet struct {
	*BaseCheevo
}

func NewMakeFirstPlanet() Cheevo {
	return &MakeFirstPlanet{
		BaseCheevo: newBaseCheevo("Make first planet"),
	}
}

func (c *MakeFirstPlanet) Init(events *twodee.GameEventHandler) {
	c.SendMessages([]string{
		"HEY",
		"HI",
		"WHAT'S UP",
	}, 2*time.Second, events)
}

func (c *MakeFirstPlanet) IsAvailable(sim *Simulation) bool {
	return true
}
