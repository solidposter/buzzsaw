package main

import (
	"log"

	"golang.org/x/net/icmp"
)

func startListener(output chan<- icmpMessage) {
	pc, err := icmp.ListenPacket("ip4:1", "0.0.0.0")
	if err != nil {
		log.Fatal(err)
	}

	for {
		rb := make([]byte, 1500)

		n, peer, err := pc.ReadFrom(rb)
		if err != nil {
			log.Fatal(err)
		}

		i := icmpMessage{
			length: n,
			peer:   peer,
			data:   rb,
		}
		// fmt.Printf("listener received ping from %+v\n", i.peer) // ugly debugging
		output <- i
	}
}
