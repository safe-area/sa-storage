package models

import "github.com/uber/h3-go"

type GetRequest struct {
	Indexes []h3.H3Index `json:"indexes"`
}

const (
	IncInfected byte = iota
	DecInfected
	IncHealthy
	DecHealthy
)

type PutRequest struct {
	Index     h3.H3Index `json:"index"`
	Timestamp int64      `json:"ts"`
	Action    byte       `json:"action"`
}
