package main

import (
	"fmt"
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

func (c *chanmonitor) start() {
	ticker := time.NewTicker(1000 * time.Millisecond)
	for {
		<-ticker.C
		fmt.Printf("dispatcher channel %v/%v\n", len(c.dispatcherChan), cap(c.dispatcherChan))
		fmt.Printf("statsengine channel %v/%v\n", len(c.statsengineChan), cap(c.statsengineChan))
		fmt.Printf("pinger channels: %v\n", len(c.pingerChans))
		for k, v := range c.pingerChans {
			if len(v) > (cap(v) / 2) {
				fmt.Printf("pinger %v channel %v/%v\n", k, len(v), cap(v))
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
