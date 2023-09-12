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
	length := []byte{0, 0}
	if _, err := io.ReadFull(cx, length); err != nil {
		return false, err
	}

	msgtype := []byte{0}
	if _, err := io.ReadFull(cx, msgtype); err != nil {
		return false, err
	}

	opcode := msgtype[0] >> 3

	keyid := msgtype[0] & 0x07

	if opcode < 1 || opcode > 9 {
		return false, nil
	}
	if keyid != 0 {
		return false, nil
	}
	return true, nil
}

var _ layer4.ConnMatcher = (*OpenvpnMatcher)(nil)
