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
	"time"
)

type statsEngine struct {
	input   chan timeReport
	targets map[string][]time.Duration
}

func newStatsEngine() *statsEngine {
	input := make(chan timeReport, 10)
	targets := make(map[string][]time.Duration)
	return &statsEngine{
		input:   input,
		targets: targets,
	}
}

func (s *statsEngine) start() {

	for {
		t := <-s.input
		rtt, exists := s.targets[t.target]
		if exists {
			rtt = append(rtt, t.rtt)
			s.targets[t.target] = rtt[1:]
			if t.rtt < 0 {
				fmt.Printf("%v %+v\n", t.target, s.targets[t.target])
			}
		} else {
			fmt.Printf("stats engine adding target %v %v\n", t.target, t.rtt)
			newList := []time.Duration{t.rtt, t.rtt, t.rtt, t.rtt, t.rtt, t.rtt, t.rtt, t.rtt, t.rtt, t.rtt}
			s.targets[t.target] = newList
		}
		// fmt.Printf("statsengne report %v %v\n", t.target, s.targets[t.target])
	}
}

func (s *statsEngine) getInput() chan timeReport {
	return s.input
}
