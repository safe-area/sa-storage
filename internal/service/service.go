package service

import (
	"github.com/safe-area/sa-storage/config"
	"github.com/safe-area/sa-storage/internal/models"
	"github.com/safe-area/sa-storage/internal/store"
	"github.com/uber/h3-go"
)

type Service interface {
	Init() error
	GetLast(indexes []h3.H3Index) map[h3.H3Index]models.HexData
	GetWithTimestamp(indexes []h3.H3Index, ts int64) map[h3.H3Index]models.HexData
	IncInfected(index h3.H3Index, ts int64) error
	DecInfected(index h3.H3Index, ts int64) error
	IncHealthy(index h3.H3Index, ts int64) error
	DecHealthy(index h3.H3Index, ts int64) error
}

func New(cfg *config.Config) Service {
	return &service{
		storage: store.New(&cfg.Storage),
		cfg:     cfg,
	}
}

type service struct {
	storage store.Storage
	cfg     *config.Config
}

func (s *service) Init() error {
	return s.storage.Init()
}

func (s *service) GetLast(indexes []h3.H3Index) map[h3.H3Index]models.HexData {
	return s.storage.GetLast(indexes)
}

func (s *service) GetWithTimestamp(indexes []h3.H3Index, ts int64) map[h3.H3Index]models.HexData {
	return s.storage.GetWithTimestamp(indexes, ts)
}

func (s *service) IncInfected(index h3.H3Index, ts int64) error {
	return s.storage.IncInfected(index, ts)
}

func (s *service) DecInfected(index h3.H3Index, ts int64) error {
	return s.storage.DecInfected(index, ts)
}

func (s *service) IncHealthy(index h3.H3Index, ts int64) error {
	return s.storage.IncHealthy(index, ts)
}

func (s *service) DecHealthy(index h3.H3Index, ts int64) error {
	return s.storage.DecHealthy(index, ts)
}
