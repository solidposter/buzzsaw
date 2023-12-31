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
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

type pinger struct {
	target string
	input  chan icmpMessage
}

func newPinger(target string) *pinger {
	input := make(chan icmpMessage, 10)
	return &pinger{
		target: target,
		input:  input,
	}
}

func (s *pinger) start(output chan timeReport) {
	var inputMessage icmpMessage
	var t0, t1 time.Time // send and receive timestamps

	dstaddr, err := net.ResolveIPAddr("ip", s.target)
	if err != nil {
		log.Fatal(err)
	}

	pc, err := icmp.ListenPacket("ip4:1", "0.0.0.0")
	if err != nil {
		log.Fatal(err)
	}
	defer pc.Close() // will never happen

	seq := 0                                          // sequence number for ICMP packet
	id := os.Getpid() & 0xffff                        // ID for icmp packets
	ticker := time.NewTicker(1000 * time.Millisecond) // one packet per second
	for {                                             // main loop
		if seq++; seq > 65534 {
			seq = 1
		}
		msg := &icmp.Message{
			Type: ipv4.ICMPTypeEcho,
			Code: 0,
			Body: &icmp.Echo{
				ID:   id,
				Seq:  seq,
				Data: []byte("buzzsaw"),
			},
		}

		wb, err := msg.Marshal(nil)
		if err != nil {
			log.Fatal(err)
		}
		// Some ugly checking during dev
		if len(s.input) != 0 {
			log.Fatalf("Input queue for %v not empty\n", s.target)
		}

		t0 = time.Now()
		if _, err := pc.WriteTo(wb, dstaddr); err != nil {
			log.Fatal(err)
		}
		//fmt.Println("pinger sent packet to", s.target)

		select {
		case <-ticker.C: // packet lost, retstart main loop
			fmt.Printf("pinger packet loss for %v\n", s.target)
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
			log.Fatal(err)
		}
		switch rm.Type {
		case ipv4.ICMPTypeEchoReply:
			//fmt.Printf("pinger received echo reply from dispatcher %v RTT %v\n", s.target, t1.Sub(t0))
			tr := timeReport{
				target: s.target,
				rtt:    t1.Sub(t0),
			}
			output <- tr
		default:
			fmt.Printf("pinger received %+v from dispatcher\n", rm.Type)
		}
		// Wait for ticker
		<-ticker.C
	}
}

func (s *pinger) getInput() chan icmpMessage {
	return s.input
}
func (s *pinger) getTarget() string {
	return s.target
}
