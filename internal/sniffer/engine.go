package sniffer

type Engine interface {
	ProcessFrame(data []byte)
	Stop() error
}
