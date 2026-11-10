package sniffer

import (
	"github.com/StewardMcCormick/tcp-analyzer/internal/packet"
	"github.com/StewardMcCormick/tcp-analyzer/internal/storage"
)

type sequentialEngine struct {
	flowStorage storage.FlowStorage
}

func NewSequentialEngine(flowStorage storage.FlowStorage) Engine {
	return &sequentialEngine{
		flowStorage: flowStorage,
	}
}

func (s *sequentialEngine) ProcessPacket(pkt packet.Packet) {
}

func (s *sequentialEngine) Stop() error {
	return nil
}
