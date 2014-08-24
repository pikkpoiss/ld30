package main

import (
	twodee "../libs/twodee"
	"math"
	"sort"
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
	for _, p := range s.Planets {
		popSum += p.GetPopulation()
		p.Update(elapsed)
		dist = p.Pos().DistanceTo(s.Sun.Pos())
		switch {
		case dist < 10:
			p.SetState(TooClose)
		case dist > 20:
			p.SetState(TooFar)
		default:
			p.SetState(Fertile)
		}
	}
	s.setPopulation(popSum)
	s.doCollisions()
}

func (s *Simulation) doCollisions() {
	var toDestroy = sort.IntSlice{}
	for index := 0; index < len(s.Planets); index++ {
		for j := index + 1; j < len(s.Planets); j++ {
			if s.Planets[index].CollidesWith(s.Planets[j]) {
				toDestroy = append(toDestroy, index)
				toDestroy = append(toDestroy, j)
			}
		}
		if s.Planets[index].CollidesWith(s.Sun) {
			toDestroy = append(toDestroy, index)
		}
		if !s.Bounds.ContainsPoint(s.Planets[index].Pos()) {
			toDestroy = append(toDestroy, index)
		}
	}

	if len(toDestroy) > 0 {
		toDestroy.Sort()
		for i := len(toDestroy) - 1; i >= 0; i-- {
			s.destroyPlanet(toDestroy[i])
		}
	}
}

func (s *Simulation) nBodyUpdate(elapsed time.Duration) {
	var dist float64
	for _, p := range s.Planets {
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

func (s *Simulation) destroyPlanet(index int) {
	var p = s.Planets[index]
	s.Planets = append(s.Planets[:index], s.Planets[index+1:]...)
	p.SetState(Exploding)
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
