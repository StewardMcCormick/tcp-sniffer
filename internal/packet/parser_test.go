package packet

import (
	"testing"
)

// Test frame constants shared by the builders below.
const (
	ethHeaderLen  = 14
	ipv4HeaderLen = 20
	tcpHeaderLen  = 20
)

// ethHeader returns a 14-byte Ethernet II header with an IPv4 EtherType.
func ethHeader() []byte {
	return []byte{
		0x00, 0x11, 0x22, 0x33, 0x44, 0x55, // Dst MAC
		0x00, 0x66, 0x77, 0x88, 0x99, 0xAA, // Src MAC
		0x08, 0x00, // EtherType IPv4
	}
}

// ipv4Header returns a 20-byte IPv4 header. fragHi and fragLo are the raw
// bytes of the "Flags + Fragment Offset" field.
func ipv4Header(fragHi, fragLo byte) []byte {
	return []byte{
		0x45, 0x00, // Version + IHL, DSCP + ECN
		0x00, 0x28, // Total Length
		0x12, 0x34, // Identification
		fragHi, fragLo, // Flags + Fragment Offset
		0x40,       // TTL
		0x06,       // Protocol (TCP)
		0x00, 0x00, // Header Checksum
		0xC0, 0xA8, 0x01, 0x01, // Src IP (192.168.1.1)
		0xC0, 0xA8, 0x01, 0x02, // Dst IP (192.168.1.2)
	}
}

// tcpHeader returns a 20-byte TCP header.
func tcpHeader() []byte {
	return []byte{
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
}

// ethIPv4TCPFrame builds a full Ethernet/IPv4/TCP frame with the given raw
// IPv4 "Flags + Fragment Offset" bytes.
func ethIPv4TCPFrame(fragHi, fragLo byte) []byte {
	frame := ethHeader()
	frame = append(frame, ipv4Header(fragHi, fragLo)...)
	frame = append(frame, tcpHeader()...)
	return frame
}

// parseFrame parses a full Ethernet/IPv4/TCP frame and fails the test on error.
func parseFrame(t *testing.T, parser *Parser, fragHi, fragLo byte) Packet {
	t.Helper()

	pkt, err := parser.Parse(ethIPv4TCPFrame(fragHi, fragLo), LinkEthernet.Key())
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}
	return pkt
}

func TestParser_Registration(t *testing.T) {
	parser := NewParser()

	keys := []Key{
		LinkEthernet.Key(),
		{KindEtherType, EtherTypeIPv4},
		{KindIPProto, IPProtoTCP},
	}

	for _, key := range keys {
		if _, ok := parser.GetHandler(key); !ok {
			t.Errorf("Handler for key %s not registered", key)
		}
	}
}

func TestParser_Parse_EthernetIPv4TCP(t *testing.T) {
	parser := NewParser()

	// DF flag set, fragment offset 0.
	pkt := parseFrame(t, parser, 0x40, 0x00)

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

	if pkt.IP.FragOffset != 0 {
		t.Errorf("Expected FragOffset 0, got %d", pkt.IP.FragOffset)
	}
	if !pkt.IP.IsFirstFragment() {
		t.Error("Expected packet to be first fragment")
	}
}

func TestParser_Parse_FragmentOffset(t *testing.T) {
	parser := NewParser()

	// The IPv4 fragment offset is a 13-bit value (max 0x1FFF == 8191);
	// the top 3 bits are flags (reserved, DF, MF) and must be masked out.
	testCases := []struct {
		name   string
		fragHi byte
		fragLo byte
		want   uint16
	}{
		{"zero offset, no flags", 0x00, 0x00, 0},
		{"min non-zero offset", 0x00, 0x01, 1},
		{"low byte only", 0x00, 0xB9, 185},
		{"boundary 0x1000", 0x10, 0x00, 4096},
		{"max offset", 0x1F, 0xFF, 8191},
		{"MF flag only, offset zero", 0x20, 0x00, 0},
		{"DF flag only, offset zero", 0x40, 0x00, 0},
		{"reserved flag only, offset zero", 0x80, 0x00, 0},
		{"all flags set, offset zero", 0xE0, 0x00, 0},
		{"MF flag plus max offset", 0x3F, 0xFF, 8191},
		{"DF flag plus max offset", 0x5F, 0xFF, 8191},
		{"all flags plus max offset", 0xFF, 0xFF, 8191},
		{"all bits set still masked to 13 bits", 0xE1, 0x2F, 0x12F},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			pkt := parseFrame(t, parser, tc.fragHi, tc.fragLo)

			if pkt.IP.FragOffset != tc.want {
				t.Errorf("FragOffset = %d, want %d", pkt.IP.FragOffset, tc.want)
			}
			if wantFirst := tc.want == 0; pkt.IP.IsFirstFragment() != wantFirst {
				t.Errorf("IsFirstFragment() = %v, want %v", pkt.IP.IsFirstFragment(), wantFirst)
			}
		})
	}
}

func TestParser_ShortPacket(t *testing.T) {
	parser := NewParser()

	full := ethIPv4TCPFrame(0x00, 0x00)

	testCases := []struct {
		name string
		data []byte
	}{
		{"shorter than ethernet header", full[:6]},
		{"one byte short of ethernet header", full[:ethHeaderLen-1]},
		{"truncated ipv4 header", full[:ethHeaderLen+10]},
		{"truncated tcp header", full[:ethHeaderLen+ipv4HeaderLen+10]},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := parser.Parse(tc.data, LinkEthernet.Key()); err == nil {
				t.Error("Expected error for short packet")
			}
		})
	}
}
