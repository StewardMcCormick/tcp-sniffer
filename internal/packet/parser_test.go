package packet

import (
	"testing"
)

func TestParser_Registration(t *testing.T) {
	parser := NewParser()

	testCases := []Key{
		LinkEthernet.Key(),
		{KindEtherType, EtherTypeIPv4},
		{KindIPProto, IPProtoTCP},
	}

	for _, tc := range testCases {
		if _, ok := parser.GetHandler(tc); !ok {
			t.Errorf("Handler for key %s not registered", tc)
		}
	}
}

func TestParser_Parse_EthernetIPv4TCP(t *testing.T) {
	parser := NewParser()

	data := []byte{
		// Ethernet header (14 bytes)
		0x00, 0x11, 0x22, 0x33, 0x44, 0x55, // Dst MAC
		0x00, 0x66, 0x77, 0x88, 0x99, 0xAA, // Src MAC
		0x08, 0x00, // EtherType IPv4

		// IPv4 header (20 bytes)
		0x45, 0x00, // Version + IHL, DSCP + ECN
		0x00, 0x28, // Total Length
		0x00, 0x00, // Identification
		0x40, 0x00, // Flags + Fragment Offset
		0x40,       // TTL
		0x06,       // Protocol (TCP)
		0x00, 0x00, // Header Checksum
		0xC0, 0xA8, 0x01, 0x01, // Src IP (192.168.1.1)
		0xC0, 0xA8, 0x01, 0x02, // Dst IP (192.168.1.2)

		// TCP header (20 bytes)
		0x30, 0x39, // Src Port (12345)
		0x00, 0x50, // Dst Port (80)
		0x00, 0x00, 0x00, 0x01, // Sequence Number
		0x00, 0x00, 0x00, 0x00, // Acknowledgement Number
		0x50,       // Data Offset + Reserved
		0x02,       // Flags (SYN)
		0x20, 0x00, // Window Size
		0x00, 0x00, // Checksum
		0x00, 0x00, // Urgent Pointer
	}

	pkt, err := parser.Parse(data, LinkEthernet.Key())
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if !pkt.HasEth {
		t.Error("Ethernet header not parsed")
	}
	if !pkt.HasIP {
		t.Error("IP header not parsed")
	}
	if !pkt.HasTCP {
		t.Error("TCP header not parsed")
	}

	if pkt.TCP.SrcPort != 12345 {
		t.Errorf("Expected SrcPort 12345, got %d", pkt.TCP.SrcPort)
	}
	if pkt.TCP.DstPort != 80 {
		t.Errorf("Expected DstPort 80, got %d", pkt.TCP.DstPort)
	}
}

func TestParser_ShortPacket(t *testing.T) {
	parser := NewParser()

	data := []byte{0x00, 0x11, 0x22, 0x33, 0x44, 0x55}

	_, err := parser.Parse(data, LinkEthernet.Key())
	if err == nil {
		t.Error("Expected error for short packet")
	}
}
