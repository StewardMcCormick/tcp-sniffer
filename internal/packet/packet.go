package packet

type Packet struct {
	Eth    EthernetLayerHeader
	IP     IPv4LayerHeader
	TCP    TCPLayerHeader
	HasEth bool
	HasIP  bool
	HasTCP bool
}

type EthernetLayerHeader struct {
	DstMAC    [6]byte
	SrcMAC    [6]byte
	EtherType uint16
}

type IPv4LayerHeader struct {
	SrcIP    [4]byte
	DstIP    [4]byte
	Protocol uint8
	IHL      uint8
}

type TCPLayerHeader struct {
	SrcPort    uint16
	DstPort    uint16
	Seq        uint32
	Ack        uint32
	Flags      uint8
	WinSize    uint16
	DataOffset uint8
}
