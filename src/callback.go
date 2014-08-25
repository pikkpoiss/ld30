package main

import (
	"time"
)

type Callback struct {
	Elapsed time.Duration
	Trigger time.Duration
	Func    func()
	Done    bool
}

func CallAfter(duration time.Duration, f func()) *Callback {
	return &Callback{
		Elapsed: 0,
		Trigger: duration,
		Func:    f,
		Done:    false,
	}
}

func (c *Callback) Update(elapsed time.Duration) {
	c.Elapsed += elapsed
	if c.Elapsed >= c.Trigger {
		c.Done = true
		c.Func()
	}
}


