package packet

import "errors"

var (
	ErrShortFrame          = errors.New("packet too short for header")
	ErrInvalidHeaderLength = errors.New("invalid header length")
	ErrUnsupportedProtocol = errors.New("unsupported protocol")
	ErrMalformedPacket     = errors.New("malformed packet")
)
