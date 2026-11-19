package sniffer

import (
	"log/slog"

	"github.com/StewardMcCormick/tcp-analyzer/internal/packet"
	"github.com/StewardMcCormick/tcp-analyzer/internal/storage"
)

type sequentialEngine struct {
	log         *slog.Logger
	flowStorage storage.FlowStorage
}

func NewSequentialEngine(log *slog.Logger, flowStorage storage.FlowStorage) Engine {
	return &sequentialEngine{
		log:         log,
		flowStorage: flowStorage,
	}
}

func (s *sequentialEngine) ProcessPacket(pkt packet.Packet) {
	if err := validatePacket(pkt); err != nil {
		s.log.Debug("Packet dropped",
			"error", err,
		)
		return
	}
}

func (s *sequentialEngine) Stop() error {
	return nil
}
