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
	for _, p := range s.Planets {
		p.GravitateToward(s.Sun.Pos())
		p.Update(elapsed)
	}
}

func (s *Simulation) AddPlanet(x, y float32) {
	s.Planets = append(s.Planets, NewPlanet(x, y))
}
