package main

import (
	"net"
	"time"
)

// see func (*icmp.PacketConn).ReadFrom(b []byte) (int, net.Addr, error)
// this struct is passed listener -> dispatcher -> pinger so the processing
// is done in the pinger.
type icmpMessage struct {
	length int
	peer   net.Addr
	data   []byte
}

// Report from pingers -> statsEngine
// a lost packet has rtt -1
type timeReport struct {
	target string
	rtt    time.Duration
}
