package sniffer

type sequentialEngine struct{}

func NewSequentialEngine() Engine {
	return &sequentialEngine{}
}

func (s *sequentialEngine) ProcessPacket(pkt *Packet) {
	// TODO implement me
	panic("implement me")
}
