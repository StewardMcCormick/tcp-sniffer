//go:build darwin

package capture

import (
	"errors"
	"log"
)

// PacketCapturer implementation for Darwin (MacOS) base on BPF
type bpfDarwinCapturer struct {
	iface string
}

func NewCapturer(iface string) (PacketCapturer, error) {
	if !isIfaceAvailable(iface) {
		return nil, errors.New("interface " + iface + " is unavailable")
	}

	return &bpfDarwinCapturer{iface: iface}, nil
}

func (b *bpfDarwinCapturer) Run() error {
	log.Printf("[INFO] Start BPF Darwin Capturer on %s", b.iface)
	return nil
}

func (b *bpfDarwinCapturer) ReadPacket() (Packet, error) {
	//TODO implement me
	panic("implement me")
}

func (b *bpfDarwinCapturer) SetBPFFilter(filter string) error {
	//TODO implement me
	panic("implement me")
}

func (b *bpfDarwinCapturer) Stop() error {
	log.Println("[INFO] Stop BPF Darwin Capturer")
	return nil
}
