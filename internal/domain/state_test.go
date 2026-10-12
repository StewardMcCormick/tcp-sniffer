package domain

import "testing"

func TestState_String(t *testing.T) {
	tests := []struct {
		name string
		s    State
		want string
	}{
		{name: "closed", s: ClosedState, want: "CLOSED"},
		{name: "established", s: EstablishedState, want: "ESTABLISHED"},
		{name: "syn_sent", s: SynSentState, want: "SYN_SENT"},
		{name: "fin_wait_1", s: FinWait1State, want: "FIN_WAIT_1"},
		{name: "fin_wait_2", s: FinWait2State, want: "FIN_WAIT_2"},
		{name: "time_wait", s: TimeWaitState, want: "TIME_WAIT"},
		{name: "syn_rcvd", s: SynRcvdState, want: "SYN_RCVD"},
		{name: "close_wait", s: CloseWaitState, want: "CLOSE_WAIT"},
		{name: "last_ack", s: LastAckState, want: "LAST_ACK"},
		{name: "invalid", s: State(100), want: "INVALID_STATE"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.String(); got != tt.want {
				t.Errorf("State(%d).String() = %q, want %q", tt.s, got, tt.want)
			}
		})
	}
}
