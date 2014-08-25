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
		if c.active.IsSatisfied(c.sim) {
			c.active.Success(c.events)
		}
		if c.active.IsReadyToDelete() {
			c.active.Delete()
			c.active = nil
		}
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
	Success(events *twodee.GameEventHandler)
	IsAvailable(sim *Simulation) bool
	IsSatisfied(sim *Simulation) bool
	IsReadyToDelete() bool
	GetLabel() string
	GetElapsed() time.Duration
	Update(elapsed time.Duration)
	Delete()
}

type BaseCheevo struct {
	label     string
	callbacks []*Callback
	elapsed   time.Duration
}

func newBaseCheevo(label string) *BaseCheevo {
	return &BaseCheevo{
		label:     label,
		callbacks: []*Callback{},
		elapsed:   0,
	}
}

func (c *BaseCheevo) GetLabel() string {
	return c.label
}

func (c *BaseCheevo) GetElapsed() time.Duration {
	return c.elapsed
}

func (c *BaseCheevo) HasPendingCallbacks() bool {
	return len(c.callbacks) > 0
}

func (c *BaseCheevo) After(d time.Duration, f func()) {
	c.callbacks = append(c.callbacks, CallAfter(d, f))
}

func (c *BaseCheevo) ClearCallbacks() {
	c.callbacks = []*Callback{}
}

func (c *BaseCheevo) Delete() {
	c.ClearCallbacks()
}

func (c *BaseCheevo) sendMessage(msg string, events *twodee.GameEventHandler) func() {
	return func() {
		events.Enqueue(NewMessageEvent(msg))
	}
}

func (c *BaseCheevo) SendMessages(messages []string, interval time.Duration, events *twodee.GameEventHandler) {
	var counter time.Duration = 0
	for i := 0; i < len(messages); i++ {
		c.After(counter, c.sendMessage(messages[i], events))
		counter += interval
	}
	c.After(counter, c.sendMessage("", events))
}

func (c *BaseCheevo) Update(elapsed time.Duration) {
	for i := len(c.callbacks) - 1; i >= 0; i-- {
		c.callbacks[i].Update(elapsed)
		if c.callbacks[i].Done {
			c.callbacks = append(c.callbacks[:i], c.callbacks[i+1:]...)
		}
	}
	c.elapsed += elapsed
}

func (c *BaseCheevo) IsReadyToDelete() bool {
	return !c.HasPendingCallbacks()
}

type MakeFirstPlanet struct {
	*BaseCheevo
	seconds    int
	hasCreated bool
	interval   time.Duration
	done       bool
	introText  []string
}

func NewMakeFirstPlanet() Cheevo {
	return &MakeFirstPlanet{
		BaseCheevo: newBaseCheevo("Make first planet"),
		hasCreated: false,
		interval:   3 * time.Second,
		done:       false,
		introText: []string{
			"",
			"HELLO",
			"WELCOME TO MY SYSTEM",
			"",
			"I SEE THAT YOU ARE ABLE TO MAKE PLANETS",
			"COULD YOU MAKE ONE FOR ME?",
			"",
			"CLICK, DRAG THE MOUSE, AND LET GO",
			"CLICK, DRAG THE MOUSE, AND LET GO",
		},
	}
}

func (c *MakeFirstPlanet) Init(events *twodee.GameEventHandler) {
	c.SendMessages(c.introText, c.interval, events)
}

func (c *MakeFirstPlanet) Success(events *twodee.GameEventHandler) {
	if !c.done {
		c.ClearCallbacks()
		c.SendMessages([]string{
			"WONDERFUL!",
		}, c.interval, events)
		c.done = true
	}
}

func (c *MakeFirstPlanet) IsAvailable(sim *Simulation) bool {
	return true
}

func (c *MakeFirstPlanet) IsSatisfied(sim *Simulation) bool {
	if len(sim.Planets) > 0 {
		c.hasCreated = true
	}

	waitTime := c.interval * time.Duration(len(c.introText)-3)
	return c.GetElapsed() > waitTime && c.hasCreated
}

//fmt.Sprintf("SEE IF YOU CAN KEEP ONE FOR LONGER THAN %v SECONDS", c.seconds),
