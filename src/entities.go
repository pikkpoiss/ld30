package main

import (
	twodee "../libs/twodee"
	"math"
	"math/rand"
	"time"
)

type PlanetaryState int32

const (
	_                  = iota
	Sun PlanetaryState = 1 << iota
	Fertile
	TooClose
	TooFar
	Exploding
	Phantom
)

var PlanetaryAnimations = map[PlanetaryState][]int{
	Sun:      []int{0, 1, 2, 3},
	Fertile:  []int{8},
	TooClose: []int{16},
	TooFar:   []int{24},
	Phantom:  []int{32},
}

type PlanetaryBody struct {
	*twodee.AnimatingEntity
	// Velocity is in units/ms.
	Velocity             twodee.Point
	Mass                 float32
	Population           float32
	MaxPopulation        float32
	PopulationGrowthRate float32
	Temperature          int32
	State                PlanetaryState
	Radius               float32
	Scale                float32
	DistToSun            float64
}

func NewSun() *PlanetaryBody {
	var (
		scale  float32 = 1.0
		length float32 = 128.0 / PxPerUnit * scale
	)
	body := &PlanetaryBody{
		AnimatingEntity: twodee.NewAnimatingEntity(
			0, 0,
			length, length,
			0,
			twodee.Step10Hz,
			[]int{0},
		),
		Mass:                 50000.0 * scale * scale,
		Population:           0.0,
		MaxPopulation:        0.0,
		PopulationGrowthRate: 0.0,
		Temperature:          27000000,
		Radius:               length / 2.0,
		Scale:                scale,
	}
	body.SetState(Sun)
	return body
}

func NewPlanet(x, y float32) *PlanetaryBody {
	var (
		scale  float32 = float32(math.Min(0.6, math.Max(0.2, rand.Float64())))
		length float32 = 128.0 / PxPerUnit * scale
	)
	body := &PlanetaryBody{
		AnimatingEntity: twodee.NewAnimatingEntity(
			x, y,
			length, length,
			0,
			twodee.Step10Hz,
			[]int{0},
		),
		Velocity:             twodee.Pt(0, 0),
		Mass:                 5000.0 * scale * scale,
		Population:           100.0,
		MaxPopulation:        0.0,
		PopulationGrowthRate: 0.0001,
		Temperature:          72,
		Radius:               length / 2.0,
		Scale:                scale,
		DistToSun:            0.0,
	}
	body.SetState(Fertile)
	body.MaxPopulation = body.Mass * 1000
	return body
}

func (p *PlanetaryBody) MoveToward(sc twodee.Point) {
	var (
		pc = p.Pos()
		dx = float64(sc.X - pc.X)
		dy = float64(sc.Y - pc.Y)
		h  = math.Hypot(dx, dy)
		vx = float32(math.Max(1, 5-h) * 0.3 * dx / h)
		vy = float32(math.Max(1, 5-h) * 0.3 * dy / h)
	)
	p.Velocity.X += (vx - p.Velocity.X)
	p.Velocity.Y += (vy - p.Velocity.Y)
}

// Calculates the PlanetaryBody's new velocity vector given an acceleration vector and time.
func (p *PlanetaryBody) CalcNewVelocity(av twodee.Point, elapsed time.Duration) {
	// Essentially, v = v0 + at
	av = av.Scale(float32(elapsed.Seconds() * 1e3))
	p.Velocity = p.Velocity.Add(av)
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
	av := twodee.Pt(float32(math.Max(1, 5-d)*0.3*avx), float32(math.Max(1, 5-d)*0.3*avy))

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
	p.Velocity.X += (fv.X - p.Velocity.X) / 30
	p.Velocity.Y += (fv.Y - p.Velocity.Y) / 30
}

func (p *PlanetaryBody) UpdatePopulation(elapsed time.Duration) {
	if p.State == Fertile {
		p.Population = p.MaxPopulation / (1 + ((p.MaxPopulation/p.Population)-1)*float32(math.Exp(-1*float64(p.PopulationGrowthRate)*float64(elapsed/time.Millisecond))))
	} else {
		p.Population = p.MaxPopulation / (1 + ((p.MaxPopulation/p.Population)-1)*float32(math.Exp(float64(p.PopulationGrowthRate)*float64(elapsed/time.Millisecond))))
	}
}

func (p *PlanetaryBody) UpdateTemperature(elapsed time.Duration) {
	if p.State == TooClose {
		p.Temperature = int32(90000.0 / math.Pow(p.DistToSun, 2.0))
	} else {
		p.Temperature = int32(5000.0 / math.Pow(p.DistToSun, 1.4))
	}

}

func (p *PlanetaryBody) Update(elapsed time.Duration) {
	p.AnimatingEntity.Update(elapsed)
	p.UpdatePopulation(elapsed)
	p.UpdateTemperature(elapsed)
	pos := p.Pos()
	ms := float32(elapsed.Seconds() * 1e3)
	dist := p.Velocity.Scale(ms)
	p.MoveTo(twodee.Pt(pos.X+dist.X, pos.Y+dist.Y))
}

func (p *PlanetaryBody) HasState(state PlanetaryState) bool {
	return p.State&state == state
}

func (p *PlanetaryBody) RemState(state PlanetaryState) {
	p.SetState(p.State & ^state)
}

func (p *PlanetaryBody) AddState(state PlanetaryState) {
	p.SetState(p.State | state)
}

func (p *PlanetaryBody) SwapState(rem, add PlanetaryState) {
	p.SetState(p.State & ^rem | add)
}

func (p *PlanetaryBody) SetState(state PlanetaryState) {
	if state != p.State {
		p.State = state
		if frames, ok := PlanetaryAnimations[p.State]; ok {
			p.SetFrames(frames)
		}
	}
}

func (p *PlanetaryBody) GetPopulation() int {
	return int(p.Population)
}

func (p *PlanetaryBody) GetTemperature() int32 {
	return int32(p.Temperature)
}

func (p *PlanetaryBody) CollidesWith(other *PlanetaryBody) bool {
	return p.Pos().DistanceTo(other.Pos()) < (p.Radius+other.Radius)*0.8
}

func (p *PlanetaryBody) SetDistToSun(dist float64) {
	p.DistToSun = dist
}
