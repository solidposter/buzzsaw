package main

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
