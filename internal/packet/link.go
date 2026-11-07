package packet

import "fmt"

type LinkKind uint16

const (
	LinkUnknown LinkKind = iota
	LinkEthernet
	LinkNull
	LinkRawIP
)

func (k LinkKind) String() string {
	switch k {
	case LinkEthernet:
		return "ethernet"
	case LinkNull:
		return "null"
	case LinkRawIP:
		return "raw_ip"
	default:
		return "unknown"
	}
}

type Key struct {
	Kind  LayerKind
	Proto uint16
}

func (k Key) String() string {
	return fmt.Sprintf("(%s,%d)", k.Kind.String(), k.Proto)
}

type LayerKind uint8

const (
	KindLink LayerKind = iota + 1
	KindEtherType
	KindIPProto
)

func (k LayerKind) String() string {
	switch k {
	case KindLink:
		return "link"
	case KindEtherType:
		return "ethertype"
	case KindIPProto:
		return "ipproto"
	default:
		return "invalid"
	}
}

func (k LinkKind) Key() Key { return Key{KindLink, uint16(k)} }
