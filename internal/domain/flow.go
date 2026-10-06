package domain

import "time"

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

type State uint8

const (
	// ---------------- General States ------------------
	ClosedState State = iota
	EstablishedState

	// ---------------- Client States -------------------
	SynSentState
	FinWait1State
	FinWait2State
	TimeWaitState

	// ---------------- Server States -------------------
	SynRcvdState
	CloseWaitState
	LastAckState
)

func (s State) String() string {
	switch s {
	case ClosedState:
		return "CLOSED"
	case EstablishedState:
		return "ESTABLISHED"
	case SynSentState:
		return "SYN_SENT"
	case FinWait1State:
		return "FIN_WAIT_1"
	case FinWait2State:
		return "FIN_WAIT_2"
	case TimeWaitState:
		return "TIME_WAIT"
	case SynRcvdState:
		return "SYN_RCVD"
	case CloseWaitState:
		return "CLOSE_WAIT"
	case LastAckState:
		return "LAST_ACK"
	default:
		return "INVALID_STATE"
	}
}

type Stats struct {
	StartTime        time.Time
	LastActivityTime time.Time

	BytesFromClient   uint64
	BytesFromServer   uint64
	PacketsFromClient uint64
	PacketsFromServer uint64

	ClientRetransmissions uint64
	ServerRetransmissions uint64

	MinRTT   float64
	MaxRTT   float64
	LastRTT  float64
	LastSRTT float64
}
