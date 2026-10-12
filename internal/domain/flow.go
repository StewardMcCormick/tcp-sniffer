package domain

type Flow struct {
	Key        FlowKey
	Server     FlowMember
	Client     FlowMember
	RolesKnown bool
}

type FlowKey struct {
	MinIp, MaxIp     [4]byte
	MinPort, MaxPort uint16
}

type FlowMember struct {
	State State
	IP    [4]byte
	Port  uint16

	HighestSeqSeen uint32
	LastAck        uint32
}
