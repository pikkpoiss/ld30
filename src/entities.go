package main

import (
	twodee "../libs/twodee"
	"math"
	"time"
)

type PlanetaryBody struct {
	*twodee.AnimatingEntity
	Velocity twodee.Point
	Mass     float32
}

func NewSun() *PlanetaryBody {
	return &PlanetaryBody{
		AnimatingEntity: twodee.NewAnimatingEntity(
			0, 0,
			32.0/PxPerUnit, 32.0/PxPerUnit,
			0,
			twodee.Step10Hz,
			[]int{
				0,
			},
		),
	}
}

func NewPlanet(x, y float32) *PlanetaryBody {
	return &PlanetaryBody{
		AnimatingEntity: twodee.NewAnimatingEntity(
			x, y,
			32.0/PxPerUnit, 32.0/PxPerUnit,
			0,
			twodee.Step10Hz,
			[]int{
				1,
			},
		),
	}
}

func (p *PlanetaryBody) MoveToward(sc twodee.Point) {
	var (
		pc = p.Pos()
		dx = float64(sc.X - pc.X)
		dy = float64(sc.Y - pc.Y)
		h  = math.Hypot(dx, dy)
		vx = float32(math.Max(1, 5-h) * 0.5 * dx / h)
		vy = float32(math.Max(1, 5-h) * 0.5 * dy / h)
	)
	p.Velocity.X += (vx - p.Velocity.X)
	p.Velocity.Y += (vy - p.Velocity.Y)
}

func (p *PlanetaryBody) Update(elapsed time.Duration) {
	p.AnimatingEntity.Update(elapsed)
	pos := p.Pos()
	p.MoveTo(twodee.Pt(pos.X+p.Velocity.X, pos.Y+p.Velocity.Y))
}
