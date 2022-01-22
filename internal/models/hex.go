package models

import "sync"

type HexCurrentProperties struct {
	Healthy    int `json:"healthy"`
	Suspicious int `json:"suspicious"`
	Infected   int `json:"infected"`
	// Users set of user Ids who was inside this hex d
	Users map[string][]TimeInterval `json:"-"`
	Mx    *sync.RWMutex
}
