package store

import (
	"github.com/nakabonne/tstorage"
	"github.com/safe-area/sa-storage/config"
	"github.com/safe-area/sa-storage/internal/models"
	"github.com/uber/h3-go"
	"path"
	"time"
)

const (
	HealthyMetric  = "healthy"
	InfectedMetric = "infected"
)

type Storage interface {
	Init() error
	GetLast(indexes []h3.H3Index) map[h3.H3Index]models.HexData
	GetWithTimestamp(indexes []h3.H3Index, ts int64) map[h3.H3Index]models.HexData
	IncInfected(index h3.H3Index, ts int64) error
	DecInfected(index h3.H3Index, ts int64) error
	IncHealthy(index h3.H3Index, ts int64) error
	DecHealthy(index h3.H3Index, ts int64) error
}

func New(cfg *config.StorageConfig) Storage {
	return &storage{
		infected: NewMetric(),
		healthy:  NewMetric(),
		cfg:      cfg,
	}
}

type storage struct {
	infected *Metric
	healthy  *Metric
	hexes    tstorage.Storage
	cfg      *config.StorageConfig
}

func (s *storage) Init() error {
	var err error
	s.hexes, err = tstorage.NewStorage(
		tstorage.WithDataPath(path.Join(s.cfg.BaseDir)),
		tstorage.WithRetention(time.Duration(s.cfg.TTL)*24*time.Hour),
	)
	return err
}

func (s *storage) GetLast(indexes []h3.H3Index) map[h3.H3Index]models.HexData {
	res := make(map[h3.H3Index]models.HexData, len(indexes))
	for _, v := range indexes {
		res[v] = models.HexData{
			Healthy:  s.healthy.Get(v),
			Infected: s.infected.Get(v),
		}
	}
	return res
}

func (s *storage) GetWithTimestamp(indexes []h3.H3Index, ts int64) map[h3.H3Index]models.HexData {
	var infected, healthy int
	var points []*tstorage.DataPoint
	var err error
	res := make(map[h3.H3Index]models.HexData, len(indexes))
	for _, v := range indexes {
		points, err = s.hexes.Select(InfectedMetric, []tstorage.Label{
			{Name: "index", Value: h3.ToString(v)},
		}, 0, ts)
		if err != nil || len(points) == 0 {
			infected = 0
		} else {
			infected = int(points[len(points)-1].Value)
		}

		points, err = s.hexes.Select(HealthyMetric, []tstorage.Label{
			{Name: "index", Value: h3.ToString(v)},
		}, 0, ts)
		if err != nil || len(points) == 0 {
			healthy = 0
		} else {
			healthy = int(points[len(points)-1].Value)
		}

		res[v] = models.HexData{
			Healthy:  healthy,
			Infected: infected,
		}
	}
	return res
}

func (s *storage) IncInfected(index h3.H3Index, ts int64) error {
	s.infected.Inc(index)
	return s.hexes.InsertRows([]tstorage.Row{
		{
			Metric: InfectedMetric,
			Labels: []tstorage.Label{
				{Name: "index", Value: h3.ToString(index)},
			},
			DataPoint: tstorage.DataPoint{Timestamp: ts, Value: float64(s.infected.Get(index))},
		},
	})
}

func (s *storage) DecInfected(index h3.H3Index, ts int64) error {
	s.infected.Dec(index)
	return s.hexes.InsertRows([]tstorage.Row{
		{
			Metric: InfectedMetric,
			Labels: []tstorage.Label{
				{Name: "index", Value: h3.ToString(index)},
			},
			DataPoint: tstorage.DataPoint{Timestamp: ts, Value: float64(s.infected.Get(index))},
		},
	})
}

func (s *storage) IncHealthy(index h3.H3Index, ts int64) error {
	s.healthy.Inc(index)
	return s.hexes.InsertRows([]tstorage.Row{
		{
			Metric: HealthyMetric,
			Labels: []tstorage.Label{
				{Name: "index", Value: h3.ToString(index)},
			},
			DataPoint: tstorage.DataPoint{Timestamp: ts, Value: float64(s.healthy.Get(index))},
		},
	})
}

func (s *storage) DecHealthy(index h3.H3Index, ts int64) error {
	s.healthy.Dec(index)
	return s.hexes.InsertRows([]tstorage.Row{
		{
			Metric: HealthyMetric,
			Labels: []tstorage.Label{
				{Name: "index", Value: h3.ToString(index)},
			},
			DataPoint: tstorage.DataPoint{Timestamp: ts, Value: float64(s.healthy.Get(index))},
		},
	})
}
