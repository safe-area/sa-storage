package service

import (
	jsoniter "github.com/json-iterator/go"
	"github.com/nats-io/nats.go"
	"github.com/safe-area/sa-storage/config"
	"github.com/safe-area/sa-storage/internal/models"
	"github.com/safe-area/sa-storage/internal/nats_provider"
	"github.com/safe-area/sa-storage/internal/store"
	"github.com/sirupsen/logrus"
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

const (
	putTemplate = "PUT_DATA_SHARD_"
	getTemplate = "GET_DATA_SHARD_"
)

func New(cfg *config.Config, provider *nats_provider.NATSProvider) Service {
	return &service{
		provider: provider,
		storage:  store.New(&cfg.Storage),
		cfg:      cfg,
	}
}

type service struct {
	storage  store.Storage
	provider *nats_provider.NATSProvider
	cfg      *config.Config
}

func (s *service) Init() error {
	var err error

	err = s.provider.Subscribe(putTemplate+s.cfg.ShardName, s.putCallback)
	if err != nil {
		return err
	}

	err = s.provider.Subscribe(getTemplate+s.cfg.ShardName, s.getCallback)
	if err != nil {
		return err
	}

	return s.storage.Init()
}

func (s *service) putCallback(msg *nats.Msg) {
	var reqs []models.PutRequest
	err := jsoniter.Unmarshal(msg.Data, &reqs)
	if err != nil {
		logrus.Errorf("PutHandler: error while unmarshalling request: %s", err)
		return
	}
	for _, req := range reqs {
		switch {
		case req.Action == models.IncInfected:
			err = s.storage.IncInfected(req.Index, req.Timestamp)
			if err != nil {
				logrus.Errorf("PutHandler: IncInfected error: %s", err)
				return
			}
		case req.Action == models.DecInfected:
			err = s.storage.DecInfected(req.Index, req.Timestamp)
			if err != nil {
				logrus.Errorf("PutHandler: DecInfected error: %s", err)
				return
			}
		case req.Action == models.IncHealthy:
			err = s.storage.IncHealthy(req.Index, req.Timestamp)
			if err != nil {
				logrus.Errorf("PutHandler: IncHealthy error: %s", err)
				return
			}
		case req.Action == models.DecHealthy:
			err = s.storage.DecHealthy(req.Index, req.Timestamp)
			if err != nil {
				logrus.Errorf("PutHandler: DecHealthy error: %s", err)
				return
			}
		default:
			logrus.Errorf("PutHandler: invalid action code: %s", err)
			return
		}
	}
	msg.Respond([]byte{})
}

func (s *service) getCallback(msg *nats.Msg) {
	var req models.GetRequest
	err := jsoniter.Unmarshal(msg.Data, &req)
	if err != nil {
		logrus.Errorf("GetHandler: error while unmarshalling request: %s", err)
		return
	}
	//todo ts field
	var ts int64
	var resp map[h3.H3Index]models.HexData
	if ts < 1 {
		resp = s.storage.GetLast(req.Indexes)
	} else {
		if err != nil {
			logrus.Errorf("GetHandler: ts query arg must be integer number or empty: %s", err)
			return
		}
		resp = s.storage.GetWithTimestamp(req.Indexes, ts)
	}
	var bs []byte
	bs, err = jsoniter.Marshal(resp)
	if err != nil {
		logrus.Errorf("GetHandler: error while unmarshalling request: %s", err)
		return
	}
	msg.Respond(bs)
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
