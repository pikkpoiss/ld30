package main

import (
	twodee "../libs/twodee"
	"fmt"
	"time"
)

type Cheevos struct {
	events  *twodee.GameEventHandler
	queue   []Cheevo
	sim     *Simulation
	active  Cheevo
	wait    time.Duration
	counter time.Duration
}

func NewCheevos(events *twodee.GameEventHandler, sim *Simulation) *Cheevos {
	return &Cheevos{
		events: events,
		queue: []Cheevo{
			NewMakeFirstPlanet(),
			NewPlanetVelocity(0.03),
			NewKeepPlanetAlive(10),
			NewMultiPlanets(2, 5),
			NewTotalPopulation(1000),
			NewMultiPlanets(3, 3),
			NewTotalPopulation(10000),
			NewMultiPlanets(4, 6),
			NewTotalPopulation(1000000),
			NewSacrifice(5000),
		},
		sim:     sim,
		active:  nil,
		wait:    5 * time.Second,
		counter: 5 * time.Second,
	}
}

func (c *Cheevos) Update(elapsed time.Duration) {
	if c.active != nil {
		c.active.Update(elapsed)
		if c.active.IsSuccess(c.sim) && !c.active.IsDone() {
			c.active.Success(c.events)
			c.active.SetDone()
		} else if c.active.IsFailure(c.sim) && !c.active.IsDone() {
			c.active.Failure(c.events)
			c.active.SetDone()
		}
		if c.active.IsReadyToDelete() {
			c.active.Delete()
			c.active = nil
			c.counter = 0
		}
	} else {
		if c.counter < c.wait {
			c.counter += elapsed
		} else {
			c.selectActive()
		}
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
	Failure(events *twodee.GameEventHandler)
	IsDone() bool
	SetDone()
	IsAvailable(sim *Simulation) bool
	IsSuccess(sim *Simulation) bool
	IsFailure(sim *Simulation) bool
	IsReadyToDelete() bool
	GetLabel() string
	GetElapsed() time.Duration
	Update(elapsed time.Duration)
	Delete()
}

type BaseCheevo struct {
	done      bool
	label     string
	callbacks []*Callback
	elapsed   time.Duration
	interval  time.Duration
	expires   time.Duration
}

func newBaseCheevo(label string, expires time.Duration) *BaseCheevo {
	return &BaseCheevo{
		label:     label,
		callbacks: []*Callback{},
		elapsed:   0,
		done:      false,
		interval:  2 * time.Second,
		expires:   expires,
	}
}

func (c *BaseCheevo) GetLabel() string {
	return c.label
}

func (c *BaseCheevo) GetElapsed() time.Duration {
	return c.elapsed
}

func (c *BaseCheevo) GetInterval() time.Duration {
	return c.interval
}

func (c *BaseCheevo) SetDone() {
	c.done = true
}

func (c *BaseCheevo) IsDone() bool {
	return c.done
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

func (c *BaseCheevo) SendMessages(messages []string, events *twodee.GameEventHandler) {
	var counter time.Duration = 0
	for i := 0; i < len(messages); i++ {
		c.After(counter, c.sendMessage(messages[i], events))
		counter += c.interval
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
	return c.done && !c.HasPendingCallbacks()
}

func (c *BaseCheevo) IsFailure(sim *Simulation) bool {
	return c.elapsed > c.expires
}

func (c *BaseCheevo) Failure(events *twodee.GameEventHandler) {
	c.ClearCallbacks()
	c.SendMessages([]string{"DARN, THAT TOOK TOO LONG"}, events)
}

// MAKE THE FIRST PLANET =======================================================

type MakeFirstPlanet struct {
	*BaseCheevo
	hasCreated bool
	introText  []string
}

func NewMakeFirstPlanet() Cheevo {
	return &MakeFirstPlanet{
		BaseCheevo: newBaseCheevo(
			"MADE YOUR FIRST PLANET",
			30*time.Second),
		hasCreated: false,
		introText: []string{
			"HELLO",
			"WELCOME TO MY SYSTEM",
			"",
			"I SEE THAT YOU ARE ABLE TO MAKE PLANETS",
			"COULD YOU MAKE ONE FOR ME?",
			"CLICK, DRAG THE MOUSE, AND LET GO",
			"CLICK, DRAG THE MOUSE, AND LET GO",
		},
	}
}

func (c *MakeFirstPlanet) Init(events *twodee.GameEventHandler) {
	c.SendMessages(c.introText, events)
}

func (c *MakeFirstPlanet) Success(events *twodee.GameEventHandler) {
	c.ClearCallbacks()
	c.SendMessages([]string{
		"WONDERFUL!",
	}, events)
}

func (c *MakeFirstPlanet) IsAvailable(sim *Simulation) bool {
	return true
}

func (c *MakeFirstPlanet) IsSuccess(sim *Simulation) bool {
	if len(sim.Planets) > 0 {
		c.hasCreated = true
	}
	waitTime := c.GetInterval() * time.Duration(len(c.introText)-1)
	return c.GetElapsed() > waitTime && c.hasCreated
}

// KEEP A PLANET ALIVE FOR MORE THAN 10 SECONDS ================================

type KeepPlanetAlive struct {
	*BaseCheevo
	introText  []string
	threshold  time.Duration
	hasPassed  bool
	planetName string
}

func NewKeepPlanetAlive(seconds int) Cheevo {
	return &KeepPlanetAlive{
		BaseCheevo: newBaseCheevo(
			fmt.Sprintf("KEPT A PLANET ALIVE FOR %v SECONDS", seconds),
			time.Duration(seconds+30)*time.Second),
		hasPassed: false,
		introText: []string{
			"HOW GOOD ARE YOU AT KEEPING ORBITS STABLE?",
			fmt.Sprintf("CAN YOU KEEP A PLANET ALIVE FOR %v SECONDS?", seconds),
		},
		threshold: time.Duration(seconds) * time.Second,
	}
}

func (c *KeepPlanetAlive) Init(events *twodee.GameEventHandler) {
	c.SendMessages(c.introText, events)
}

func (c *KeepPlanetAlive) Success(events *twodee.GameEventHandler) {
	c.ClearCallbacks()
	c.SendMessages([]string{
		fmt.Sprintf("GREAT! %v HAS LIVED A LONG TIME", c.planetName),
	}, events)
}

func (c *KeepPlanetAlive) IsAvailable(sim *Simulation) bool {
	for _, p := range sim.Planets {
		if p.Age > c.threshold {
			return false
		}
	}
	return true
}

func (c *KeepPlanetAlive) IsSuccess(sim *Simulation) bool {
	if c.hasPassed == false {
		for _, p := range sim.Planets {
			if p.Age > c.threshold {
				c.hasPassed = true
				c.planetName = p.Name
				break
			}
		}
	}
	waitTime := c.GetInterval() * time.Duration(len(c.introText))
	return c.GetElapsed() > waitTime && c.hasPassed
}

// MAKE A PLANET'S VELOCITY X ==================================================

type PlanetVelocity struct {
	*BaseCheevo
	introText  []string
	velocity   float32
	hasPassed  bool
	planetName string
}

func NewPlanetVelocity(velocity float32) Cheevo {
	return &PlanetVelocity{
		BaseCheevo: newBaseCheevo(
			"MADE A PLANET GO FAST",
			15*time.Second),
		hasPassed: false,
		introText: []string{
			"CAN YOU MAKE A PLANET GO FAST?",
			"CLICK, DRAG THE MOUSE A DISTANCE, AND LET GO",
		},
		velocity: velocity,
	}
}

func (c *PlanetVelocity) Init(events *twodee.GameEventHandler) {
	c.SendMessages(c.introText, events)
}

func (c *PlanetVelocity) Success(events *twodee.GameEventHandler) {
	c.ClearCallbacks()
	c.SendMessages([]string{
		fmt.Sprintf("WOAH! %v IS A SPEEDY ONE", c.planetName),
	}, events)
}

func (c *PlanetVelocity) IsAvailable(sim *Simulation) bool {
	return true
}

func (c *PlanetVelocity) IsSuccess(sim *Simulation) bool {
	if c.hasPassed == false {
		pt := twodee.Pt(0, 0)
		for _, p := range sim.Planets {
			if p.Velocity.DistanceTo(pt) >= c.velocity {
				c.hasPassed = true
				c.planetName = p.Name
				break
			}
		}
	}
	waitTime := c.GetInterval() * time.Duration(len(c.introText)-1)
	return c.GetElapsed() > waitTime && c.hasPassed
}

// N PLANETS FOR Y SECONDS =====================================================

type MultiPlanets struct {
	*BaseCheevo
	introText   []string
	planetCount int32
	threshold   time.Duration
	hasPassed   bool
}

func NewMultiPlanets(count int32, seconds int32) Cheevo {
	return &MultiPlanets{
		BaseCheevo: newBaseCheevo(
			fmt.Sprintf("HAD %v PLANETS LIVE FOR %v SECONDS EACH", count, seconds),
			40*time.Second),
		hasPassed: false,
		introText: []string{
			fmt.Sprintf("I WANT TO SEE %v PLANETS AT ONCE", count),
			fmt.Sprintf("HAVE THEM SURVIVE FOR %v SECONDS EACH", seconds),
		},
		planetCount: count,
		threshold:   time.Duration(seconds) * time.Second,
	}
}

func (c *MultiPlanets) Init(events *twodee.GameEventHandler) {
	c.SendMessages(c.introText, events)
}

func (c *MultiPlanets) Success(events *twodee.GameEventHandler) {
	c.ClearCallbacks()
	c.SendMessages([]string{
		"SO MANY PLANETS!",
	}, events)
}

func (c *MultiPlanets) IsAvailable(sim *Simulation) bool {
	return len(sim.Planets) < int(c.planetCount)
}

func (c *MultiPlanets) IsSuccess(sim *Simulation) bool {
	if c.hasPassed == false {
		var count int32 = 0
		for _, p := range sim.Planets {
			if p.Age > c.threshold {
				count += 1
			}
		}
		if count >= c.planetCount {
			c.hasPassed = true
		}
	}
	waitTime := c.GetInterval() * time.Duration(len(c.introText))
	return c.GetElapsed() > waitTime && c.hasPassed
}

// TOTAL POPULATION ============================================================

type TotalPopulation struct {
	*BaseCheevo
	introText  []string
	population int32
	hasPassed  bool
}

func NewTotalPopulation(population int32) Cheevo {
	return &TotalPopulation{
		BaseCheevo: newBaseCheevo(
			fmt.Sprintf("ACHIEVED %v POPULATION", population),
			300*time.Second),
		hasPassed: false,
		introText: []string{
			"I YEARN FOR MORE LIFE",
			fmt.Sprintf("PRODUCE %v TOTAL SOULS", population),
		},
		population: population,
	}
}

func (c *TotalPopulation) Init(events *twodee.GameEventHandler) {
	c.SendMessages(c.introText, events)
}

func (c *TotalPopulation) Success(events *twodee.GameEventHandler) {
	c.ClearCallbacks()
	c.SendMessages([]string{
		fmt.Sprintf("I FEEL THE WARMTH OF %v TINY BODIES", c.population),
	}, events)
}

func (c *TotalPopulation) IsAvailable(sim *Simulation) bool {
	return sim.GetPopulation() < int(c.population)
}

func (c *TotalPopulation) IsSuccess(sim *Simulation) bool {
	if c.hasPassed == false {
		if sim.GetPopulation() >= int(c.population) {
			c.hasPassed = true
		}
	}
	waitTime := c.GetInterval() * time.Duration(len(c.introText))
	return c.GetElapsed() > waitTime && c.hasPassed
}

// SACRIFICE ===================================================================

type Sacrifice struct {
	*BaseCheevo
	hasPassed  bool
	hasFailed  bool
	population int32
	target     *PlanetaryBody
	planetName string
	events     *twodee.GameEventHandler
	obsFire    int
	obsColl    int
}

func NewSacrifice(population int32) Cheevo {
	return &Sacrifice{
		BaseCheevo: newBaseCheevo(
			"MADE THE ULTIMATE SACRIFICE",
			300*time.Second),
		hasPassed:  false,
		hasFailed:  false,
		population: population,
	}
}

func (c *Sacrifice) Init(events *twodee.GameEventHandler) {
	c.SendMessages([]string{
		"I AM UNFULFILLED",
		"YOU HAVE BROUGHT SO MANY SOULS TO ME",
		fmt.Sprintf("BRING %v INTO MY GREATNESS", c.planetName),
	}, events)
	c.obsFire = events.AddObserver(PlanetFireDeath, c.OnFireDeath)
	c.obsColl = events.AddObserver(PlanetCollision, c.OnCollision)
	c.events = events
}

func (c *Sacrifice) OnFireDeath(evt twodee.GETyper) {
	switch event := evt.(type) {
	case *PlanetEvent:
		if event.Planet == c.target {
			c.hasPassed = true
		}
	}
}

func (c *Sacrifice) OnCollision(evt twodee.GETyper) {
	switch event := evt.(type) {
	case *PlanetEvent:
		if event.Planet == c.target {
			c.hasFailed = true
		}
	}
}

func (c *Sacrifice) Failure(events *twodee.GameEventHandler) {
	c.ClearCallbacks()
	c.SendMessages([]string{
		"I AM AN ANGRY SOL!",
	}, events)
}

func (c *Sacrifice) Success(events *twodee.GameEventHandler) {
	c.ClearCallbacks()
	c.SendMessages([]string{
		fmt.Sprintf("%v IS MINE!", c.planetName),
	}, events)
}

func (c *Sacrifice) IsAvailable(sim *Simulation) bool {
	for _, p := range sim.Planets {
		if int(p.Population) > int(c.population) {
			c.target = p
			c.planetName = p.Name
			return true
		}
	}
	return false
}

func (c *Sacrifice) IsSuccess(sim *Simulation) bool {
	waitTime := c.GetInterval() * time.Duration(3)
	return c.GetElapsed() > waitTime && c.hasPassed
}

func (c *Sacrifice) IsFailure(sim *Simulation) bool {
	return c.BaseCheevo.IsFailure(sim) || c.hasFailed
}

func (c *Sacrifice) Delete() {
	c.BaseCheevo.Delete()
	c.events.RemoveObserver(PlanetFireDeath, c.obsFire)
	c.events.RemoveObserver(PlanetCollision, c.obsColl)
}
