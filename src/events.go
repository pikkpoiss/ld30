package main

import (
	twodee "../libs/twodee"
)

const (
	GameIsClosing twodee.GameEventType = iota
	sentinel
)

const (
	NumGameEventTypes = int(sentinel)
)
