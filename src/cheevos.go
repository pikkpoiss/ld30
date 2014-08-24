package main

import (
	twodee "../libs/twodee"
	"time"
)

type Cheevos struct {
	Events *twodee.GameEventHandler
}

func NewCheevos(events *twodee.GameEventHandler) *Cheevos {
	return nil
}

func (c *Cheevos) Update(elapsed time.Duration) {
}

func (c *Cheevos) Delete() {
}
