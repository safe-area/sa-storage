package models

import "sync"

type AtomicInt struct {
	v  int
	mx *sync.RWMutex
}

func NewAtomicInt() *AtomicInt {
	return &AtomicInt{
		v:  0,
		mx: &sync.RWMutex{},
	}
}

func (i *AtomicInt) Add(delta int) {
	i.mx.Lock()
	i.v += delta
	i.mx.Unlock()
}

func (i *AtomicInt) Get() int {
	i.mx.RLock()
	v := i.v
	i.mx.RUnlock()
	return v
}
