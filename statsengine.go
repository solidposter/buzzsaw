package main

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
