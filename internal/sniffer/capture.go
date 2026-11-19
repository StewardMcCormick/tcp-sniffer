package sniffer

import (
	"errors"
	"net"
)

var (
	ErrCannotInitializeSniffer = errors.New("cannot initialize error")
	ErrCannotCloseSniffer      = errors.New("cannot stop sniffer")
)

func isIfaceAvailable(iface string) bool {
	_, err := net.InterfaceByName(iface)

	return err == nil
}
