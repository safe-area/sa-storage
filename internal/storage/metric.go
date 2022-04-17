package store

import (
	"github.com/safe-area/sa-storage/internal/models"
	"github.com/uber/h3-go"
	"sync"
)

type Metric struct {
	m  map[h3.H3Index]*models.AtomicInt
	mx *sync.RWMutex
}

func NewMetric() *Metric {
	return &Metric{
		m:  make(map[h3.H3Index]*models.AtomicInt),
		mx: &sync.RWMutex{},
	}
}

func (m *Metric) Inc(index h3.H3Index) {
	if _, ok := m.m[index]; !ok {
		m.mx.Lock()
		// double check for decrease mutex locks
		if _, ok = m.m[index]; !ok {
			m.m[index] = models.NewAtomicInt()
		}
		m.mx.Unlock()
	}
	m.mx.RLock()
	m.m[index].Add(1)
	m.mx.RUnlock()
}

func (m *Metric) Dec(index h3.H3Index) {
	if _, ok := m.m[index]; !ok {
		m.mx.Lock()
		// double check for decrease mutex locks
		if _, ok = m.m[index]; !ok {
			m.m[index] = models.NewAtomicInt()
		}
		m.mx.Unlock()
	}
	m.mx.RLock()
	m.m[index].Add(-1)
	m.mx.RUnlock()
}

func (m *Metric) Get(index h3.H3Index) int {
	if _, ok := m.m[index]; !ok {
		return 0
	}
	m.mx.RLock()
	v := m.m[index].Get()
	m.mx.RUnlock()
	return v
}
