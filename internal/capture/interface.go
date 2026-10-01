package capture

import "time"

type Packet struct {
	Data      []byte
	Timestamp time.Time
}

type PacketCapturer interface {
	Run(iface string) error

	ReadPacket() (Packet, error)

	SetBPFFilter(filter string) error

	Stop() error
}
