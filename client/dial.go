package client

import (
	"context"
	"net"
	"time"
)

// dialContext provides connection dialing with timeout/keepalive.
type dialContext struct {
	Timeout   time.Duration
	KeepAlive time.Duration
}

// DialContext implements the dialer.
func (d *dialContext) DialContext(ctx context.Context, network, address string) (net.Conn, error) {
	dialer := &net.Dialer{
		Timeout:   d.Timeout,
		KeepAlive: d.KeepAlive,
	}
	return dialer.DialContext(ctx, network, address)
}
