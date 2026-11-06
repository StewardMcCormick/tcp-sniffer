package packet

type Packet struct {
	IP  IPv4LayerHeader
	TCP TCPLayerHeader
}

type EthernetLayerHeader struct {
	SrcMac    uint8
	DstMac    uint8
	EtherType uint8
}

type IPv4LayerHeader struct {
	SrcIP    [4]byte
	DstIP    [4]byte
	Protocol uint8
}

type TCPLayerHeader struct {
	SrcPort uint16
	DstPort uint16
	Seq     uint32
	Ack     uint32
	Flags   uint16
	WinSize uint16
}
