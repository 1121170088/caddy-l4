package l4openvpntcp

import (
	"github.com/caddyserver/caddy/v2"
	"github.com/mholt/caddy-l4/layer4"
	"io"
)

func init() {
	caddy.RegisterModule(OpenvpnMatcher{})
}
type OpenvpnMatcher struct {
}

func (OpenvpnMatcher) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "layer4.matchers.openvpntcp",
		New: func() caddy.Module { return new(OpenvpnMatcher) },
	}
}

func (m OpenvpnMatcher) Match(cx *layer4.Connection) (bool, error) {
	// read pkg length
	lengthBs := []byte{0, 0}
	if _, err := io.ReadFull(cx, lengthBs); err != nil {
		return false, err
	}
	length := int(lengthBs[1]) | (int(lengthBs[0])<<8)
	msgtype := []byte{0}
	if _, err := io.ReadFull(cx, msgtype); err != nil {
		return false, err
	}

	opcode := msgtype[0] >> 3

	keyid := msgtype[0] & 0x07
	if keyid != 0 {
		return false, nil
	}
	if length < 14 {
		return false, nil
	}
	switch opcode {
	case 1:
		return false, nil
	case 2:
		return false, nil
	case 3:
		return false, nil
	case 4:
		if length > 2000 {
			return false, nil
		}
	case 5:
		if length != 22 {
			return false, nil
		}
	case 6:
		return false, nil
	case 7:
		if length != 14 {
			return false, nil
		}
	case 8:
		if length != 26 {
			return false, nil
		}
	case 9:
		if length > 2000 {
			return false, nil
		}
	default:
		return false, nil
	}

	return true, nil
}

var _ layer4.ConnMatcher = (*OpenvpnMatcher)(nil)
