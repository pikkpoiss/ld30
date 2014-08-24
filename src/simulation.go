package main

import (
	"math"
	"time"
)

const (
	// 6.67e-11 m^3kg^-1s^-2
	GravitationalConst = 6.67e-11
	// Play GC in m^3kg^-1ms^-2
	GravConst = 5e-8
)

type Simulation struct {
	Sun                 *PlanetaryBody
	Planets             []*PlanetaryBody
	AggregatePopulation int
}

func NewSimulation() *Simulation {
	return &Simulation{
		Sun:                 NewSun(),
		Planets:             []*PlanetaryBody{},
		AggregatePopulation: 0,
	}
}

func (s *Simulation) Update(elapsed time.Duration) {
	s.NBodyUpdate(elapsed)
	var popSum = 0
	for _, p := range s.Planets {
		popSum += p.GetPopulation()
		p.Update(elapsed)
		dist = float64(p.Pos().DistanceTo(s.Sun.Pos()))
		switch {
		case dist < 10:
			p.SetState(TooClose)
		case dist > 15:
			p.SetState(TooFar)
		default:
			p.SetState(Fertile)
		}
	}
	s.SetPopulation(popSum)
}

func (s *Simulation) NBodyUpdate(elapsed time.Duration) {
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

func (s *Simulation) AddPlanet(x, y float32) {
	s.Planets = append(s.Planets, NewPlanet(x, y))
}

func (s *Simulation) SetPopulation(population int) {
	s.AggregatePopulation = population
}

func (s *Simulation) GetPopulation() int {
	return s.AggregatePopulation
}
