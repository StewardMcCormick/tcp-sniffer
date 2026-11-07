package sniffer

import (
	"fmt"

	"github.com/StewardMcCormick/tcp-analyzer/internal/packet"
)

// Константы DLT из https://www.tcpdump.org/linktypes.html
const (
	dltNull     uint32 = 0
	dltEthernet uint32 = 1
	dltRaw      uint32 = 12
)

func LinkKindFromDLT(dlt uint32) (packet.LinkKind, error) {
	switch dlt {
	case dltEthernet:
		return packet.LinkEthernet, nil
	case dltNull:
		return packet.LinkNull, nil
	case dltRaw:
		return packet.LinkRawIP, nil
	default:
		return packet.LinkUnknown, fmt.Errorf("unsupported DLT: %d", dlt)
	}
}
