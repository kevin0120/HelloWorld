package utils

import (
	"sync"
)

const (
	DEFAULT_SEQUENCE = 1
	MAX_SEQUENCE     = 99999
)

func CreateSequence(maxSequence uint32) *Sequence {
	if maxSequence == 0 {
		maxSequence = MAX_SEQUENCE
	}
	return &Sequence{
		mtx:         sync.Mutex{},
		sequence:    DEFAULT_SEQUENCE,
		maxSequence: maxSequence,
	}
}

type Sequence struct {
	sequence    uint32
	mtx         sync.Mutex
	maxSequence uint32
}

//func (c *Sequence) setSequence(i uint32) {
//	c.mtx.Lock()
//	defer c.mtx.Unlock()
//	if c.sequence == c.maxSequence {
//		c.sequence = DEFAULT_SEQUENCE
//		return
//	}
//	c.sequence = i
//}

func (c *Sequence) ResetSequence() {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	c.sequence = DEFAULT_SEQUENCE
}

func (c *Sequence) IncSequence() {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	if c.sequence == c.maxSequence {
		c.sequence = DEFAULT_SEQUENCE
		return
	}
	c.sequence += 1
}

func (c *Sequence) ValidSequence(seq uint32) bool {
	if seq < DEFAULT_SEQUENCE {
		return false
	}
	if seq > MAX_SEQUENCE {
		return false
	}
	return true
}

func (c *Sequence) GetSequenceWithInc() uint32 {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	seq := c.sequence
	if c.sequence == c.maxSequence {
		c.sequence = DEFAULT_SEQUENCE
	} else {
		c.sequence += 1
	}
	return seq
}

func (c *Sequence) GetSequence() uint32 {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	return c.sequence
}
