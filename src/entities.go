package main

import (
	twodee "../libs/twodee"
)

func NewSun() *twodee.AnimatingEntity {
	return twodee.NewAnimatingEntity(
		0, 0,
		32.0/PxPerUnit, 32.0/PxPerUnit,
		0,
		twodee.Step10Hz,
		[]int{
			0,
		},
	)
}

func NewPlanet(x, y float32) *twodee.AnimatingEntity {
	return twodee.NewAnimatingEntity(
		x, y,
		32.0/PxPerUnit, 32.0/PxPerUnit,
		0,
		twodee.Step10Hz,
		[]int{
			1,
		},
	)
}

