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
)

type dispatcher struct {
	pingers map[string]chan icmpMessage
	input   chan icmpMessage
}

func newDispatcher() *dispatcher {
	pingers := make(map[string]chan icmpMessage)
	input := make(chan icmpMessage, 10)
	return &dispatcher{
		pingers: pingers,
		input:   input,
	}
}

func (d *dispatcher) start() {
	for {
		packet := <-d.input
		target := packet.peer.String()
		output, exists := d.pingers[target]
		if exists {
			// fmt.Println("dispatcher received from listener", packet.peer.String())
			output <- packet
		} else {
			fmt.Printf("dispatcher unknown peer %+v\n", packet.peer)
		}
	}
}

func (d *dispatcher) getInput() chan icmpMessage {
	return d.input
}

func (d *dispatcher) addPinger(target string, clientchannel chan icmpMessage) {
	_, exists := d.pingers[target]
	if exists {
		fmt.Println("target already exists", target)
		return
	} else {
		fmt.Println("dispatcher adding target", target)
		d.pingers[target] = clientchannel
	}
}
