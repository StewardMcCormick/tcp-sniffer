package packet

import (
	"fmt"
)

type Handler func(data []byte, pkt *Packet) (consumed int, next Key, err error)

type Parser struct {
	handlers map[Key]Handler
}

func NewParser() *Parser {
	p := &Parser{
		handlers: make(map[Key]Handler, 8),
	}
	p.registerDefaultHandlers()
	return p
}

func (p *Parser) registerDefaultHandlers() {
	// Link layer handlers
	p.Register(LinkEthernet.Key(), parseEthernet)
	p.Register(LinkNull.Key(), parseNull)
	p.Register(LinkRawIP.Key(), parseRawIP)

	// EtherType handlers
	p.Register(Key{KindEtherType, EtherTypeIPv4}, parseIPv4)

	// IP protocol handlers
	p.Register(Key{KindIPProto, IPProtoTCP}, parseTCP)
}

func (p *Parser) Register(k Key, h Handler) {
	p.handlers[k] = h
}

func (p *Parser) GetHandler(k Key) (Handler, bool) {
	h, ok := p.handlers[k]
	return h, ok
}

const maxLayers = 8

func (p *Parser) Parse(data []byte, root Key) (Packet, error) {
	var pkt Packet
	key, off := root, 0

	for depth := 0; depth < maxLayers && off < len(data); depth++ {
		h, ok := p.handlers[key]
		if !ok {
			break
		}

		n, next, err := h(data[off:], &pkt)
		if err != nil {
			return pkt, fmt.Errorf("layer %d (%s): %w", depth, key, err)
		}
		if n < 0 || off+n > len(data) {
			return pkt, fmt.Errorf("layer %d (%s): %w", depth, key, ErrInvalidHeaderLength)
		}

		off += n
		if next == (Key{}) {
			break
		}
		if next == key {
			return pkt, fmt.Errorf("layer %d (%s): handler returned same key", depth, key)
		}
		key = next
	}

	return pkt, nil
}
