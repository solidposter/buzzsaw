package main

//
// Copyright (c) 2023 Tony Sarendal <tony@polarcap.org>
//
// Permission to use, copy, modify, and distribute this software for any
// purpose with or without fee is hereby granted, provided that the above
// copyright notice and this permission notice appear in all copies.
//
// THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES
// WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF
// MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR
// ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES
// WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN
// ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF
// OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.
//

import (
	"log/slog"
	"net"
	"os"
	"time"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

type pinger struct {
	target     string
	resolvedIP net.IPAddr
	input      chan icmpMessage
}

func newPinger(target string) (*pinger, error) {
	resolvedIP, err := net.ResolveIPAddr("ip", target)
	if err != nil {
		return nil, err
	}
	input := make(chan icmpMessage, 10)
	return &pinger{
		target:     target,
		resolvedIP: *resolvedIP,
		input:      input,
	}, nil
}

func (s *pinger) start(output chan timeReport) {
	var inputMessage icmpMessage
	var t0, t1 time.Time // send and receive timestamps

	dstaddr, err := net.ResolveIPAddr("ip", s.target)
	if err != nil {
		slog.Error("ResolveIPAddr() failed", "error", err)
		os.Exit(1)
	}

	pc, err := icmp.ListenPacket("ip4:1", "0.0.0.0")
	if err != nil {
		slog.Error("ListenPacket failed", "error", err)
		os.Exit(1)
	}
	defer pc.Close() // will never happen

	seq := 0                                          // sequence number for ICMP packet
	id := os.Getpid() & 0xffff                        // ID for icmp packets
	ticker := time.NewTicker(1000 * time.Millisecond) // one packet per second
	for {                                             // main loop
		if seq++; seq > 65534 {
			seq = 1
		}
		sendEcho(pc, dstaddr, seq, id)
		t0 = time.Now()

		// Some ugly checking during dev
		if len(s.input) != 0 {
			slog.Error("Input queu should be empty", "target", s.target)
			os.Exit(1)
		}

		select {
		case <-ticker.C: // packet lost, retstart main loop
			slog.Debug("Packet loss detected", "peer", s.target)
			tr := timeReport{
				target: s.target,
				rtt:    -1,
			}
			output <- tr
			continue
		case inputMessage = <-s.input: // packet received
			t1 = time.Now()
		}

		rm, err := icmp.ParseMessage(1, inputMessage.data[:inputMessage.length])
		if err != nil {
			slog.Error("ParseMessage failed", "error", err)
			os.Exit(1)
		}
		switch rm.Type {
		case ipv4.ICMPTypeEchoReply:
			slog.Debug("Reply recived", "peer", s.target, "rtt", t1.Sub(t0))
			tr := timeReport{
				target: s.target,
				rtt:    t1.Sub(t0),
			}
			output <- tr
		default:
			slog.Warn("Non-reply reveived", "peer", s.target, "type", rm.Type)
		}
		// Wait for ticker
		<-ticker.C
	}
}

func (s *pinger) getInput() chan icmpMessage {
	return s.input
}
func (s *pinger) getTarget() string {
	return s.resolvedIP.String()
}

func sendEcho(conn *icmp.PacketConn, target *net.IPAddr, seq int, id int) {
	icmpBody := newIcmpEcho(seq, id)

	encoded, err := icmpBody.Marshal(nil)
	if err != nil {
		slog.Error("Marshall of ICMP failed", "error", err)
		os.Exit(1)
	}

	if _, err := conn.WriteTo(encoded, target); err != nil {
		slog.Error("WriteTo() failed", "error", err)
		os.Exit(1)
	}
}

func newIcmpEcho(seq int, id int) *icmp.Message {
	return &icmp.Message{
		Type: ipv4.ICMPTypeEcho,
		Code: 0,
		Body: &icmp.Echo{
			ID:   id,
			Seq:  seq,
			Data: []byte("buzzsaw"),
		},
	}
}
