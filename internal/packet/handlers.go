package packet

import (
	"encoding/binary"
)

func parseEthernet(data []byte, pkt *Packet) (int, Key, error) {
	const ethHeaderLen = 14
	if len(data) < ethHeaderLen {
		return 0, Key{}, ErrShortFrame
	}

	copy(pkt.Eth.DstMAC[:], data[0:6])
	copy(pkt.Eth.SrcMAC[:], data[6:12])
	pkt.Eth.EtherType = binary.BigEndian.Uint16(data[12:14])
	pkt.HasEth = true

	return ethHeaderLen, Key{KindEtherType, pkt.Eth.EtherType}, nil
}

func parseNull(data []byte, pkt *Packet) (int, Key, error) {
	const nullHeaderLen = 4
	if len(data) < nullHeaderLen {
		return 0, Key{}, ErrShortFrame
	}

	af := binary.NativeEndian.Uint32(data[0:4])

	if af == 2 {
		return nullHeaderLen, Key{KindEtherType, EtherTypeIPv4}, nil
	}

	return nullHeaderLen, Key{}, nil
}

func parseRawIP(data []byte, pkt *Packet) (int, Key, error) {
	if len(data) < 1 {
		return 0, Key{}, ErrShortFrame
	}

	version := data[0] >> 4

	if version == 4 {
		return 0, Key{KindEtherType, EtherTypeIPv4}, nil
	}

	return 0, Key{}, nil
}

func parseIPv4(data []byte, pkt *Packet) (int, Key, error) {
	const minIPv4HeaderLen = 20
	if len(data) < minIPv4HeaderLen {
		return 0, Key{}, ErrShortFrame
	}

	ihl := int(data[0]&0x0F) * 4
	if ihl < minIPv4HeaderLen || len(data) < ihl {
		return 0, Key{}, ErrInvalidHeaderLength
	}

	copy(pkt.IP.SrcIP[:], data[12:16])
	copy(pkt.IP.DstIP[:], data[16:20])
	pkt.IP.Protocol = data[9]
	pkt.IP.IHL = uint8(ihl)
	pkt.HasIP = true

	return ihl, Key{KindIPProto, uint16(pkt.IP.Protocol)}, nil
}

func parseTCP(data []byte, pkt *Packet) (int, Key, error) {
	const minTCPHeaderLen = 20
	if len(data) < minTCPHeaderLen {
		return 0, Key{}, ErrShortFrame
	}

	pkt.TCP.SrcPort = binary.BigEndian.Uint16(data[0:2])
	pkt.TCP.DstPort = binary.BigEndian.Uint16(data[2:4])
	pkt.TCP.Seq = binary.BigEndian.Uint32(data[4:8])
	pkt.TCP.Ack = binary.BigEndian.Uint32(data[8:12])
	pkt.TCP.DataOffset = (data[12] >> 4) * 4
	pkt.TCP.Flags = data[13]
	pkt.TCP.WinSize = binary.BigEndian.Uint16(data[14:16])
	pkt.HasTCP = true

	tcpHeaderLen := int(pkt.TCP.DataOffset)
	if tcpHeaderLen < minTCPHeaderLen || len(data) < tcpHeaderLen {
		return 0, Key{}, ErrInvalidHeaderLength
	}

	return tcpHeaderLen, Key{}, nil
}
