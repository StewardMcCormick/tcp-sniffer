package sniffer

import (
	"net"
	"time"
)

type Config struct {
	Interface      string        `yaml:"interface" env-required:"true"`
	SnapshotLength uint16        `yaml:"max_snapshot_length" env-default:"65536"`
	Promisc        bool          `yaml:"promiscuous_mode" env-default:"false"`
	Timeout        time.Duration `yaml:"timeout" env-default:"25ms"`
}

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

	return err == nil
}
