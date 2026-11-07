package sniffer

import (
	"errors"
	"fmt"
	"log"
	"net"

	"github.com/StewardMcCormick/tcp-analyzer/internal/config"
	"github.com/StewardMcCormick/tcp-analyzer/internal/packet"
	"github.com/google/gopacket/pcap"
)

// PacketCapturer implementation using gopacket
type pcapCapturer struct {
	cfg    config.SnifferConfig
	engine Engine
	parser *packet.Parser
}

func NewPcapCapturer(cfg config.SnifferConfig) (PacketCapturer, error) {
	if !isIfaceAvailable(cfg.Interface) {
		return nil, errors.New("interface " + cfg.Interface + " is unavailable")
	}

	var engine Engine
	switch cfg.Engine.Type {
	case config.EngineTypeSequential:
		engine = NewSequentialEngine()
	default:
		return nil, errors.New("unsupported engine type")
	}

	parser := packet.NewParser()

	return &pcapCapturer{
		cfg:    cfg,
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

	log.Printf("[INFO] Starting capture on %s, link type: %s", b.cfg.Interface, rootKind)

	for {
		data, ci, err := h.ReadPacketData()
		if err != nil {
			var netErr net.Error
			if ok := errors.Is(err, netErr); (ok && netErr.Timeout()) || errors.Is(err, pcap.NextErrorTimeoutExpired) {
				log.Println("[INFO] No packets for capturing.")
				continue
			}
			log.Printf("[ERROR] Packet capturing error: %v\n", err)
			continue
		}

		if ci.CaptureLength < ci.Length {
			data = data[:ci.CaptureLength]
		}

		pkt, err := b.parser.Parse(data, root)
		if err != nil {
			log.Printf("[WARNING] Packet parsing error: %v\n", err)
			continue
		}

		b.engine.ProcessPacket(pkt)
	}
}

func (b *pcapCapturer) Stop() error {
	if err := b.engine.Stop(); err != nil {
		return fmt.Errorf("%w - %w", ErrCannotCloseSniffer, err)
	}

	log.Println("[INFO] Gopacket Capturer stopped")
	return nil
}
