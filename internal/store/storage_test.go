package store

import (
	"github.com/safe-area/sa-storage/config"
	"github.com/stretchr/testify/assert"
	"github.com/uber/h3-go"
	"testing"
)

func TestStorage(t *testing.T) {
	s := New(&config.StorageConfig{
		TTL:     1,
		BaseDir: "data",
	})
	if err := s.Init(); err != nil {
		t.Fatal(err)
	}
	if err := s.IncHealthy(1, 1); err != nil {
		t.Fatal(err)
	}
	if err := s.IncHealthy(1, 2); err != nil {
		t.Fatal(err)
	}
	if err := s.DecHealthy(1, 3); err != nil {
		t.Fatal(err)
	}
	if err := s.IncHealthy(1, 4); err != nil {
		t.Fatal(err)
	}
	if err := s.DecInfected(1, 1); err != nil {
		t.Fatal(err)
	}
	if err := s.DecInfected(1, 2); err != nil {
		t.Fatal(err)
	}
	if err := s.IncInfected(1, 3); err != nil {
		t.Fatal(err)
	}
	if err := s.IncInfected(1, 4); err != nil {
		t.Fatal(err)
	}
	if err := s.IncInfected(1, 5); err != nil {
		t.Fatal(err)
	}
	gl := s.GetLast([]h3.H3Index{1})
	assert.Equal(t, 2, gl[1].Healthy)
	assert.Equal(t, 1, gl[1].Infected)

	gwt := s.GetWithTimestamp([]h3.H3Index{1}, 6)
	assert.Equal(t, 2, gwt[1].Healthy)
	assert.Equal(t, 1, gwt[1].Infected)

	gwt = s.GetWithTimestamp([]h3.H3Index{1}, 3)
	assert.Equal(t, 2, gwt[1].Healthy)
	assert.Equal(t, -2, gwt[1].Infected)
}
