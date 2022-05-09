package api

import (
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"github.com/safe-area/sa-storage/internal/models"
	"github.com/safe-area/sa-storage/internal/service"
	"github.com/sirupsen/logrus"
	"github.com/uber/h3-go"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttprouter"
	"strconv"
	"time"
)

type Server struct {
	r    *fasthttprouter.Router
	serv *fasthttp.Server
	svc  service.Service
	port string
}

func New(svc service.Service, port string) *Server {
	innerRouter := fasthttprouter.New()
	innerHandler := innerRouter.Handler
	s := &Server{
		innerRouter,
		&fasthttp.Server{
			ReadTimeout:  time.Duration(10) * time.Second,
			WriteTimeout: time.Duration(10) * time.Second,
			IdleTimeout:  time.Duration(10) * time.Second,
			Handler:      innerHandler,
		},
		svc,
		port,
	}

	s.r.GET("/api/v1/get", s.GetHandler)
	s.r.POST("/api/v1/put", s.PutHandler)

	return s
}

func (s *Server) GetHandler(ctx *fasthttp.RequestCtx, ps fasthttprouter.Params) {
	body := ctx.PostBody()
	var req models.GetRequest
	err := jsoniter.Unmarshal(body, &req)
	if err != nil {
		logrus.Errorf("GetHandler: error while unmarshalling request: %s", err)
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}
	tsStr := string(ctx.QueryArgs().Peek("ts"))
	var resp map[h3.H3Index]models.HexData
	if tsStr == "" {
		resp = s.svc.GetLast(req.Indexes)
	} else {
		var ts int64
		ts, err = strconv.ParseInt(tsStr, 10, 64)
		if err != nil {
			logrus.Errorf("GetHandler: ts query arg must be integer number or empty: %s", err)
			ctx.SetStatusCode(fasthttp.StatusBadRequest)
			return
		}
		resp = s.svc.GetWithTimestamp(req.Indexes, ts)
	}
	var bs []byte
	bs, err = jsoniter.Marshal(resp)
	if err != nil {
		logrus.Errorf("GetHandler: error while unmarshalling request: %s", err)
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}
	ctx.Write(bs)
	ctx.SetContentType("application/json")
	ctx.SetStatusCode(fasthttp.StatusOK)
}

func (s *Server) PutHandler(ctx *fasthttp.RequestCtx, ps fasthttprouter.Params) {
	body := ctx.PostBody()
	var reqs []models.PutRequest
	err := jsoniter.Unmarshal(body, &reqs)
	if err != nil {
		logrus.Errorf("PutHandler: error while unmarshalling request: %s", err)
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}
	for _, req := range reqs {
		switch {
		case req.Action == models.IncInfected:
			err = s.svc.IncInfected(req.Index, req.Timestamp)
			if err != nil {
				logrus.Errorf("PutHandler: IncInfected error: %s", err)
				ctx.SetStatusCode(fasthttp.StatusBadRequest)
				return
			}
		case req.Action == models.DecInfected:
			err = s.svc.DecInfected(req.Index, req.Timestamp)
			if err != nil {
				logrus.Errorf("PutHandler: DecInfected error: %s", err)
				ctx.SetStatusCode(fasthttp.StatusBadRequest)
				return
			}
		case req.Action == models.IncHealthy:
			err = s.svc.IncHealthy(req.Index, req.Timestamp)
			if err != nil {
				logrus.Errorf("PutHandler: IncHealthy error: %s", err)
				ctx.SetStatusCode(fasthttp.StatusBadRequest)
				return
			}
		case req.Action == models.DecHealthy:
			err = s.svc.DecHealthy(req.Index, req.Timestamp)
			if err != nil {
				logrus.Errorf("PutHandler: DecHealthy error: %s", err)
				ctx.SetStatusCode(fasthttp.StatusBadRequest)
				return
			}
		default:
			logrus.Errorf("PutHandler: invalid action code: %s", err)
			ctx.SetStatusCode(fasthttp.StatusBadRequest)
			return
		}
	}
	ctx.SetStatusCode(fasthttp.StatusOK)
}

func (s *Server) Start() error {
	return fmt.Errorf("server start: %s", s.serv.ListenAndServe(s.port))
}
func (s *Server) Shutdown() error {
	return s.serv.Shutdown()
}
