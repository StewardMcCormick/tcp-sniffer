package sniffer

import (
	"errors"
	"net"
)

var (
	ErrCannotInitializeSniffer = errors.New("cannot initialize error")
	ErrCannotCloseSniffer      = errors.New("cannot stop sniffer")
)

type PacketCapturer interface {
	Run() error

	Stop() error
}

func isIfaceAvailable(iface string) bool {
	_, err := net.InterfaceByName(iface)

	return err == nil
}
