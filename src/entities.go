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
