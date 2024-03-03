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
	"time"
)

type chanmonitor struct {
	dispatcherChan  chan icmpMessage
	pingerChans     map[string]chan icmpMessage
	statsengineChan chan timeReport
}

func newChanmonitor() *chanmonitor {
	p := make(map[string]chan icmpMessage, 10)
	return &chanmonitor{
		pingerChans: p,
	}
}

const channelWarning string = "Channel overload"

func (c *chanmonitor) start() {
	slog.Info("Starting channel monitor")
	ticker := time.NewTicker(1000 * time.Millisecond)
	for {
		<-ticker.C
		if len(c.dispatcherChan) > cap(c.dispatcherChan)/2 {
			slog.Warn(channelWarning, "channel", "dispatcher", "length", len(c.dispatcherChan), "capacity", cap(c.dispatcherChan))
		}

		if len(c.statsengineChan) > cap(c.statsengineChan)/2 {
			slog.Warn(channelWarning, "channel", "statsengine", "length", len(c.statsengineChan), "capacity", cap(c.statsengineChan))
		}

		//fmt.Printf("pinger channels: %v\n", len(c.pingerChans))
		for k, v := range c.pingerChans {
			if len(v) > cap(v)/2 {
				slog.Warn(channelWarning, "channel", "pinger", "pinger", k, "length", len(v), "capacity", cap(v))
			}
		}
	}
}

func (c *chanmonitor) addDispatcher(d chan icmpMessage) {
	c.dispatcherChan = d
}

func (c *chanmonitor) addPinger(target string, d chan icmpMessage) {
	c.pingerChans[target] = d
}

func (c *chanmonitor) addStatsengine(d chan timeReport) {
	c.statsengineChan = d
}
