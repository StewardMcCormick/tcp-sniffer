package sniffer

import "github.com/StewardMcCormick/tcp-analyzer/internal/packet"

type Engine interface {
	ProcessPacket(pkt packet.Packet)
	Stop() error
}
