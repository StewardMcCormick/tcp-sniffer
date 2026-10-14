package sniffer

import (
	"errors"
	"log"
)

// PacketCapturer implementation using gopacket
type pcapCapturer struct {
	iface  string
	engine Engine
}

func NewPcapCapturer(iface string, engine Engine) (PacketCapturer, error) {
	if !isIfaceAvailable(iface) {
		return nil, errors.New("interface " + iface + " is unavailable")
	}

	return &pcapCapturer{iface: iface, engine: engine}, nil
}

func (b *pcapCapturer) Run() error {
	return nil
}

func (b *pcapCapturer) ReadPacket() (Packet, error) {
	// TODO implement me
	panic("implement me")
}

func (b *pcapCapturer) SetBPFFilter(filter string) error {
	// TODO implement me
	panic("implement me")
}

func (b *pcapCapturer) Stop() error {
	log.Println("[INFO] Stop BPF Darwin Capturer")
	return nil
}
