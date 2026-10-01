package capture

// PacketCapturer implementation for Darwin (MacOS) base on BPF

type bpfDarwinCapturer struct {
}

func (b bpfDarwinCapturer) Run(iface string) error {
	//TODO implement me
	panic("implement me")
}

func (b bpfDarwinCapturer) ReadPacket() (Packet, error) {
	//TODO implement me
	panic("implement me")
}

func (b bpfDarwinCapturer) SetBPFFilter(filter string) error {
	//TODO implement me
	panic("implement me")
}

func (b bpfDarwinCapturer) Stop() error {
	//TODO implement me
	panic("implement me")
}
