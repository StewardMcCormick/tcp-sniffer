package domain

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
