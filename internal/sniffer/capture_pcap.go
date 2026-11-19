package sniffer

import (
	"errors"
	"fmt"
	"log/slog"
	"net"

	"github.com/StewardMcCormick/tcp-analyzer/internal/config"
	"github.com/StewardMcCormick/tcp-analyzer/internal/packet"
	"github.com/StewardMcCormick/tcp-analyzer/internal/storage"
	"github.com/google/gopacket/pcap"
)

type Engine interface {
	ProcessPacket(pkt packet.Packet)
	Stop() error
}

type pcapCapturer struct {
	cfg    config.SnifferConfig
	log    *slog.Logger
	engine Engine
	parser *packet.Parser
}

func NewPcapCapturer(cfg config.SnifferConfig, log *slog.Logger, storage storage.FlowStorage) (*pcapCapturer, error) {
	if !isIfaceAvailable(cfg.Interface) {
		return nil, errors.New("interface " + cfg.Interface + " is unavailable")
	}

	var engine Engine
	switch cfg.Engine.Type {
	case config.EngineTypeSequential:
		engine = NewSequentialEngine(log, storage)
	default:
		return nil, errors.New("unsupported engine type")
	}

	parser := packet.NewParser()

	return &pcapCapturer{
		cfg:    cfg,
		log:    log,
		engine: engine,
		parser: parser,
	}, nil
}

func (b *pcapCapturer) Run() error {
	h, err := pcap.OpenLive(b.cfg.Interface, int32(b.cfg.SnapshotLength), b.cfg.Promisc, b.cfg.Timeout)
	if err != nil {
		return fmt.Errorf("%w - %w", ErrCannotInitializeSniffer, err)
	}
	defer h.Close()

	rootKind, err := LinkKindFromDLT(uint32(h.LinkType()))
	if err != nil {
		return fmt.Errorf("unsupported link type %v: %w", h.LinkType(), err)
	}
	root := rootKind.Key()

	b.log.Info("Starting capture",
		"interface", b.cfg.Interface,
		"link_type", rootKind,
	)

	for {
		data, ci, err := h.ReadPacketData()
		if err != nil {
			var netErr net.Error
			if ok := errors.Is(err, netErr); (ok && netErr.Timeout()) || errors.Is(err, pcap.NextErrorTimeoutExpired) {
				continue
			}
			b.log.Error("Packet capturing error",
				"error", err,
			)
			continue
		}

		if ci.CaptureLength < ci.Length {
			data = data[:ci.CaptureLength]
		}

		pkt, err := b.parser.Parse(data, root)
		if err != nil {
			b.log.Warn("Cannot parse packet",
				"error", err,
			)
			continue
		}

		b.engine.ProcessPacket(pkt)
	}
}

func (b *pcapCapturer) Stop() error {
	if err := b.engine.Stop(); err != nil {
		return fmt.Errorf("%w - %w", ErrCannotCloseSniffer, err)
	}

	b.log.Info("Gopacket Capturer stopped")
	return nil
}
