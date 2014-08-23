package main

import (
	"time"
)

type Simulation struct {
	Sun     *PlanetaryBody
	Planets []*PlanetaryBody
}

func NewSimulation() *Simulation {
	return &Simulation{
		Sun:     NewSun(),
		Planets: []*PlanetaryBody{},
	}
}

func (s *Simulation) Update(elapsed time.Duration) {
	var centroid = s.Sun.Pos().Scale(s.Sun.Mass)
	var weight = s.Sun.Mass
	for _, p := range s.Planets {
		centroid = centroid.Add(p.Pos().Scale(p.Mass))
		weight += p.Mass
	}
	centroid = centroid.Scale(1.0 / weight)
	for _, p := range s.Planets {
		p.GravitateToward(centroid)
		p.Update(elapsed)
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
}

func (s *Simulation) AddPlanet(x, y float32) {
	s.Planets = append(s.Planets, NewPlanet(x, y))
}
