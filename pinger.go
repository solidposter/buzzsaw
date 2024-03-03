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

const IcmpDataString = "buzzsaw"

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
		tSend := time.Now()

		for { // Run waitloop until ticker.C
			select {
			case <-ticker.C:
				slog.Debug("Packet loss detected", "peer", s.target)
				tr := timeReport{
					target: s.target,
					rtt:    -1,
				}
				output <- tr
			case inputMessage = <-s.input:
				slog.Debug("Message recived", "peer", s.target)
				if !isValidResponse(inputMessage, seq, id) {
					slog.Warn("Invalid packet", "peer", s.target)
					continue
				}

				tRecv := time.Now() // Valid packet received
				tr := timeReport{
					target: s.target,
					rtt:    tRecv.Sub(tSend),
				}
				output <- tr
				<-ticker.C
			}
			break
		}
	}
}

func (s *pinger) getInput() chan icmpMessage {
	return s.input
}
func (s *pinger) getTarget() string {
	return s.resolvedIP.String()
}

func isValidResponse(message icmpMessage, seq int, id int) bool {
	icmpMessage, err := icmp.ParseMessage(1, message.data[:message.length])
	if err != nil {
		slog.Error("ParseMessage failed", "error", err)
		os.Exit(1)
	}
	if icmpMessage.Type != ipv4.ICMPTypeEchoReply {
		return false
	}
	echoReply, ok := icmpMessage.Body.(*icmp.Echo)
	if !ok {
		slog.Error("Failed to decode ICMP echo reply")
		os.Exit(1)
	}

	if seq == echoReply.Seq && id == echoReply.ID && IcmpDataString == string(echoReply.Data) {
		return true
	} else {
		return false
	}
}

func sendEcho(conn *icmp.PacketConn, target *net.IPAddr, seq int, id int) {
	icmpMessage := newIcmpEcho(seq, id)

	encoded, err := icmpMessage.Marshal(nil)
	if err != nil {
		slog.Error("Marshall of ICMP failed", "error", err)
		os.Exit(1)
	}

	if _, err := conn.WriteTo(encoded, target); err != nil {
		slog.Error("Sending packet failed", "error", err)
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
			Data: []byte(IcmpDataString),
		},
	}
}
