package sniffer

import (
	"errors"
	"fmt"
	"log"
	"net"

	"github.com/StewardMcCormick/tcp-analyzer/internal/config"
	"github.com/google/gopacket/pcap"
)

// PacketCapturer implementation using gopacket
type pcapCapturer struct {
	cfg    config.SnifferConfig
	engine Engine
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

	return &pcapCapturer{
		cfg:    cfg,
		engine: engine,
	}, nil
}

func (b *pcapCapturer) Run() error {
	h, err := pcap.OpenLive(b.cfg.Interface, int32(b.cfg.SnapshotLength), b.cfg.Promisc, b.cfg.Timeout)
	if err != nil {
		return fmt.Errorf("%w - %w", ErrCannotInitializeSniffer, err)
	}
	defer h.Close()

	for {
		data, ci, err := h.ReadPacketData()
		if err != nil {
			var netErr net.Error
			if ok := errors.Is(err, netErr); (ok && netErr.Timeout()) || errors.Is(err, pcap.NextErrorTimeoutExpired) {
				log.Println("[INFO] No any packets for capturing.")
				continue
			}
			log.Printf("[ERROR] Packet capturing error: %v\n", err)
			continue
		}

		log.Printf("[INFO] New packet: %v, timestamp: %v\n, Total length: %d", data, ci.Timestamp, ci.Length) // TODO: will delete in future commits
		log.Printf("[INFO] Captured data length: %d", len(data))

		if err != nil {
			log.Printf("[WARNING] Packet parsing error: %v\n", err)
		}

		b.engine.ProcessFrame(data)
	}
}

func (b *pcapCapturer) Stop() error {
	if err := b.engine.Stop(); err != nil {
		return fmt.Errorf("%w - %w", ErrCannotCloseSniffer, err)
	}

	log.Println("[INFO] Gopacket Capturer stopped")
	return nil
}
