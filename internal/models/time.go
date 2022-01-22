package models

// TimeInterval [TimeIn;TimeOut] in Unix
type TimeInterval struct {
	In     int64
	Out    int64
	UserId string
}
