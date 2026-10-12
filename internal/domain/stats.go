package domain

import "time"

type Stats struct {
	StartTime        time.Time
	LastActivityTime time.Time

	BytesFromClient   uint64
	BytesFromServer   uint64
	PacketsFromClient uint64
	PacketsFromServer uint64

	ClientRetransmissions uint64
	ServerRetransmissions uint64

	MinRTT   float64
	MaxRTT   float64
	LastRTT  float64
	LastSRTT float64
}
