package main

import (
	"time"
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
	var centroid = s.Sun.Pos().Scale(s.Sun.Mass)
	var weight = s.Sun.Mass
	var PopulationAggregator = 0
	for _, p := range s.Planets {
		centroid = centroid.Add(p.Pos().Scale(p.Mass))
		weight += p.Mass
	}
	centroid = centroid.Scale(1.0 / weight)
	for _, p := range s.Planets {
		p.GravitateToward(centroid)
		p.Update(elapsed)
		PopulationAggregator += p.GetPopulation()
		dist := p.Pos().DistanceTo(s.Sun.Pos())
		switch {
		case dist < 10:
			p.SetState(TooClose)
		case dist > 15:
			p.SetState(TooFar)
		default:
			p.SetState(Fertile)
		}
	}
	s.SetPopulation(PopulationAggregator)
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
