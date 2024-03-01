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
	"log/slog"
)

type dispatcher struct {
	pingers map[string]chan icmpMessage
	input   chan icmpMessage
}

func newDispatcher(qlen int) *dispatcher {
	pingers := make(map[string]chan icmpMessage)
	input := make(chan icmpMessage, qlen)
	return &dispatcher{
		pingers: pingers,
		input:   input,
	}
}

func (d *dispatcher) start() {
	slog.Info("Starting dispatcher")
	for {
		packet := <-d.input
		target := packet.peer.String()
		slog.Debug("Dispatcher received packet", "peer", target)
		output, exists := d.pingers[target]
		if exists {
			output <- packet
		} else {
			slog.Warn("Dispatcher received packet from unknown peer", "peer", target)
		}
	}
}

func (d *dispatcher) getInput() chan icmpMessage {
	return d.input
}

func (d *dispatcher) addPinger(target string, clientchannel chan icmpMessage) error {
	_, exists := d.pingers[target]
	if exists {
		err := fmt.Errorf("duplicate target %s", target)
		return err
	} else {
		d.pingers[target] = clientchannel
		slog.Debug("Pinger added to dispatcher", "target", target)
		return nil
	}
}
