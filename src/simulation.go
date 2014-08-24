package main

import (
	twodee "../libs/twodee"
	"math"
	"time"
)

const (
	// 6.67e-11 m^3kg^-1s^-2
	GravitationalConst = 6.67e-11
	// Play GC in m^3kg^-1ms^-2
	GravConst    = 5e-8
	BoundsBuffer = 10.0
)

type Simulation struct {
	Sun                 *PlanetaryBody
	Planets             []*PlanetaryBody
	AggregatePopulation int
	MaxPopulation       int
	Events              *twodee.GameEventHandler
	Bounds              twodee.Rectangle
}

func NewSimulation(bounds twodee.Rectangle, events *twodee.GameEventHandler) *Simulation {
	return &Simulation{
		Sun:                 NewSun(),
		Planets:             []*PlanetaryBody{},
		AggregatePopulation: 0,
		MaxPopulation:       0,
		Events:              events,
		Bounds: twodee.Rect(
			bounds.Min.X-BoundsBuffer,
			bounds.Min.Y-BoundsBuffer,
			bounds.Max.X+BoundsBuffer,
			bounds.Max.Y+BoundsBuffer,
		),
	}
}

func (s *Simulation) Update(elapsed time.Duration) {
	var (
		dist   float32
		popSum = 0
	)
	s.nBodyUpdate(elapsed)
	s.Sun.Update(elapsed)
	for _, p := range s.Planets {
		popSum += p.GetPopulation()
		p.Update(elapsed)
		if p.HasState(Dying) || p.HasState(Dead) {
			continue
		}
		dist = p.Pos().DistanceTo(s.Sun.Pos())
		p.SetDistToSun(float64(dist))
		switch {
		case dist < 12:
			p.SetState(TooClose)
		case dist > 30:
			p.SetState(TooFar)
		default:
			p.SetState(Fertile)
		}
	}
	s.setPopulation(popSum)
	s.doCollisions()
	s.doRemoveDeadPlanets()
}

func (s *Simulation) doCollisions() {
	for index := 0; index < len(s.Planets); index++ {
		for j := index + 1; j < len(s.Planets); j++ {
			if s.Planets[index].CollidesWith(s.Planets[j]) {
				s.destroyPlanet(index, Colliding)
				s.destroyPlanet(j, Colliding)
			}
		}
		if s.Planets[index].CollidesWith(s.Sun) {
			s.destroyPlanet(index, Exploding)
		}
		if !s.Bounds.ContainsPoint(s.Planets[index].Pos()) {
			s.Planets[index].SetState(Dead)
		}
	}
}

func (s *Simulation) doRemoveDeadPlanets() {
	for i := len(s.Planets) - 1; i >= 0; i-- {
		if s.Planets[i].HasState(Dead) {
			s.removePlanet(i)
		}
	}
}

func (s *Simulation) nBodyUpdate(elapsed time.Duration) {
	var dist float64
	for _, p := range s.Planets {
		if p.HasState(Dying) || p.HasState(Dead) {
			continue
		}
		// First, we must handle the sun...
		dist = float64(s.Sun.Pos().DistanceTo(p.Pos()))
		var force = s.Sun.Pos().Sub(p.Pos()).Scale(s.Sun.Mass * p.Mass).Scale(float32(math.Pow(dist, -3)))
		for _, p2 := range s.Planets {
			if p == p2 {
				continue
			}
			dist = float64(p2.Pos().DistanceTo(p.Pos()))
			force = force.Add(p2.Pos().Sub(p.Pos()).Scale(p2.Mass * p.Mass).Scale(float32(math.Pow(dist, -3))))
		}
		// Don't forget the gravitational constant and p's mass.
		var accel = force.Scale(GravConst).Scale(1 / p.Mass)
		p.CalcNewVelocity(accel, elapsed)
	}
}

func (s *Simulation) AddPlanet(p *PlanetaryBody) {
	s.Planets = append(s.Planets, p)
}

func (s *Simulation) removePlanet(index int) {
	s.Planets = append(s.Planets[:index], s.Planets[index+1:]...)
}

func (s *Simulation) destroyPlanet(index int, state PlanetaryState) {
	s.Planets[index].Destroy(state)
}

func (s *Simulation) setPopulation(population int) {
	s.AggregatePopulation = population
	if population > s.MaxPopulation {
		s.MaxPopulation = population
	}
}

func (s *Simulation) GetPopulation() int {
	return s.AggregatePopulation
}

func (s *Simulation) GetMaxPopulation() int {
	return s.MaxPopulation
}
