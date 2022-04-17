package store

import (
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

func TestMetric(t *testing.T) {
	m := NewMetric()
	wg := &sync.WaitGroup{}
	for i := 0; i < 3000; i++ {
		wg.Add(1)
		go func(m *Metric, i int) {
			if i%3 == 0 {
				m.Dec(1)
				m.Inc(2)
			} else {
				m.Inc(1)
				m.Dec(2)
			}
			wg.Done()
		}(m, i)
	}
	wg.Wait()
	assert.Equal(t, 1000, m.Get(1))
	assert.Equal(t, -1000, m.Get(2))
}
