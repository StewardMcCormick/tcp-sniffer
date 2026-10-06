package capture

import (
	"net"
	"time"
)

type Packet struct {
	Data      []byte
	Timestamp time.Time
}

type PacketCapturer interface {
	Run() error

	ReadPacket() (Packet, error)

	SetBPFFilter(filter string) error

	Stop() error
}

func isIfaceAvailable(iface string) bool {
	_, err := net.InterfaceByName(iface)
	if err != nil {
		return false
	}

	return true
}
