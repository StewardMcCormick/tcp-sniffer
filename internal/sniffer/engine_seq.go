package sniffer

import (
	"time"
)

type sequentialEngine struct{}

func NewSequentialEngine() Engine {
	return &sequentialEngine{}
}

func (s *sequentialEngine) ProcessFrame(data []byte) {
	time.Sleep(500*time.Millisecond) //nolint TODO: will fix in future commits
}

func (s *sequentialEngine) Stop() error  {
	return nil
}
