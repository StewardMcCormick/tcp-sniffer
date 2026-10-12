package sniffer

type Engine interface {
	ProcessPacket(pkt *Packet)
}
