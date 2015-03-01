package proxy

import (
	"io"

	"github.com/mcuadros/exmongodb/protocol"
)

// proxyMessage proxies a message, possibly it's response, and possibly a
// follow up call.
func (p *Proxy) Handle(h *protocol.MsgHeader, c io.ReadWriter, s io.ReadWriter) error {
	if err := h.WriteTo(s); err != nil {
		p.Log.Error(err)
		return err
	}

	if _, err := io.CopyN(s, c, int64(h.MessageLength-headerLen)); err != nil {
		p.Log.Error(err)
		return err
	}

	// For Ops with responses we proxy the raw response message over.
	if h.OpCode.HasResponse() {
		//stats.BumpSum(p.stats, "message.with.response", 1)
		if err := protocol.CopyMessage(c, s); err != nil {
			p.Log.Error(err)
			return err
		}
	}

	return nil
}
