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

func (p *PlanetaryBody) GravitateToward(sc twodee.Point) {
	var (
		pc  = p.Pos()
		avx = float64(sc.X - pc.X)
		avy = float64(sc.Y - pc.Y)
		d   = math.Hypot(avx, avy)
	)
	// Normalize vector and include sensible constraints.
	avx = avx / d
	avy = avy / d
	av := twodee.Pt(float32(math.Max(1, 5-d)*0.2*avx), float32(math.Max(1, 5-d)*0.2*avy))

	// There are two possible orthogonal 'circulation' vectors.
	cv1 := twodee.Pt(-av.Y, av.X)
	cv2 := twodee.Pt(av.Y, -av.X)
	cv := cv1

	// Compute whichever circulation vector is closer to our present vector.
	// cos(theta) = A -dot- B / ||A||*||B||
	dp1 := p.Velocity.X*cv1.X + p.Velocity.Y*cv1.Y
	denom := math.Sqrt(float64(p.Velocity.X*p.Velocity.X + p.Velocity.Y*p.Velocity.Y))
	theta1 := dp1 / float32(denom)
	dp2 := p.Velocity.X*cv2.X + p.Velocity.Y*cv2.Y
	theta2 := dp2 / float32(denom)
	if theta1 >= theta2 {
		cv = cv1
	} else {
		cv = cv2
	}

	// Now do some vector addition.
	fv := twodee.Pt(av.X+cv.X, av.Y+cv.Y)
	p.Velocity.X += (fv.X - p.Velocity.X)
	p.Velocity.Y += (fv.Y - p.Velocity.Y)
}

func (p *PlanetaryBody) Update(elapsed time.Duration) {
	p.AnimatingEntity.Update(elapsed)
	pos := p.Pos()
	p.MoveTo(twodee.Pt(pos.X+p.Velocity.X, pos.Y+p.Velocity.Y))
}
