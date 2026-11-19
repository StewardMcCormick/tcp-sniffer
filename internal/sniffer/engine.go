package sniffer

import (
	"errors"

	"github.com/StewardMcCormick/tcp-analyzer/internal/packet"
)

var ErrPacketIsNotFirstFragment = errors.New("packet is not first packet")

func validatePacket(pkt packet.Packet) error {
	if !pkt.IP.IsFirstFragment() {
		return ErrPacketIsNotFirstFragment
	}

	return nil
}
