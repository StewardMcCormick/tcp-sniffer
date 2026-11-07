package sniffer

import (
	"time"

	"github.com/StewardMcCormick/tcp-analyzer/internal/packet"
)

type sequentialEngine struct{}

func NewSequentialEngine() Engine {
	return &sequentialEngine{}
}

func (s *sequentialEngine) ProcessPacket(pkt packet.Packet) {
	time.Sleep(500*time.Millisecond) //nolint TODO: will fix in future commits
}

func (s *sequentialEngine) Stop() error  {
	return nil
}
