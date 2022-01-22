package storage

import (
	"github.com/safe-area/sa-storage/internal/models"
	"github.com/uber/h3-go"
)

type Storage interface {
	Put()
}

type storage struct {
	current map[h3.H3Index]*models.HexLoad
}
