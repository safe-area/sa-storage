package store

import (
	"github.com/safe-area/sa-storage/internal/models"
	"github.com/uber/h3-go"
	"sync"
)

type Storage interface {
}

func New() Storage {
	return &storage{
		Hexes: make(map[h3.H3Index]*models.HexCurrentProperties),
	}
}

type storage struct {
	Hexes   map[h3.H3Index]*models.HexCurrentProperties
	HexesMx *sync.RWMutex
	Users   map[string]*models.User
	UsersMx *sync.RWMutex
}

func (s *storage) GetAll(indexes []h3.H3Index) map[h3.H3Index]*models.HexCurrentProperties {
	res := make(map[h3.H3Index]*models.HexCurrentProperties, len(indexes))
	s.HexesMx.RLock()
	for _, v := range indexes {
		res[v] = s.Hexes[v]
	}
	s.HexesMx.RUnlock()
	return res
}
