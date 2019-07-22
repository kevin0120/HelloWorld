// Copyright © 2019 Bjørn Erik Pedersen <bjorn.erik.pedersen@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

// Package debounce provides a debouncer func. The most typical use case would be
// the user typing a text into a form; the UI needs an update, but let's wait for
// a break.
package throttle

import (
	"sync"
	"time"
)

// New returns a debounced function that takes another functions as its argument.
// This function will be called when the debounced function stops being called
// for the given duration.
// The debounced function can be invoked with different functions, if needed,
// the last one will win.
func New(after time.Duration) func(f func()) {
	d := &throttler{after: after,
		ticker : time.NewTicker(after),
	}
	go func() {
		for{
			<-d.ticker.C
			if d.have_func{
				d.func1()
				d.have_func=false
			}
		}
	}()
	return func(f func()) {
		d.add(f)
	}
}

type throttler struct {
	mu    sync.Mutex
	after time.Duration
	ticker *time.Ticker
	have_func bool
	func1 func()
}

func (d *throttler) add(f func()) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.func1=f
	d.have_func=true
}
