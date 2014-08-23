package main

import (
	twodee "../libs/twodee"
)

const (
	GameIsClosing twodee.GameEventType = iota
	PlayBackgroundMusic
	sentinel
)

const (
	NumGameEventTypes = int(sentinel)
)
