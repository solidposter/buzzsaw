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
