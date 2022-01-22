package models

import "github.com/uber/h3-go"

const (
	Healthy    = UserState(0)
	Suspicious = UserState(1)
	Infected   = UserState(2)
)

type UserState byte

type User struct {
	Id        string
	Hex       h3.H3Index
	State     UserState
	LastState UserState
}
