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

type statsEngine struct {
	input chan timeReport
}

func newStatsEngine(qlen int) *statsEngine {
	input := make(chan timeReport, qlen)
	return &statsEngine{
		input: input,
	}
}

func (s *statsEngine) start() {
	slog.Info("Starting statsengine")
	targets := make(map[string][]time.Duration)
	ticker := time.NewTicker(1 * time.Second)
	for {
		select {
		case t := <-s.input:
			slog.Debug("timeReport received", "target", t.target, "rtt", t.rtt)
			rtt, exists := targets[t.target]
			if exists {
				rtt = append(rtt, t.rtt)
				targets[t.target] = rtt[1:]
			} else {
				slog.Debug("Target added to statsengine", "target", t.target)
				newList := []time.Duration{t.rtt, t.rtt, t.rtt, t.rtt, t.rtt, t.rtt, t.rtt, t.rtt, t.rtt, t.rtt}
				targets[t.target] = newList
				if t.rtt == -1 {
					slog.Info("Target status is down", "target", t.target)
				}
			}
		case <-ticker.C:
			logStatusToDown(targets)
			logStatusToUp(targets)
		}
	}
}

func (s *statsEngine) getInput() chan timeReport {
	return s.input
}

func logStatusToDown(targets map[string][]time.Duration) {
	for target, rttList := range targets {
		if isStatusToDown(rttList) {
			slog.Info("Status changed to down", "target", target)
		}
	}
}

// A change from valid RTT from three consecutive -1 means means status change to down
func isStatusToDown(rttList []time.Duration) bool {
	if l := len(rttList); l > 4 {
		if rttList[l-4] != -1 {
			if rttList[l-1] == -1 && rttList[l-2] == -1 && rttList[l-3] == -1 {
				return true
			}
		}
	}
	return false
}

func logStatusToUp(targets map[string][]time.Duration) {
	for target, rttList := range targets {
		if isStatusToUp(rttList) {
			slog.Info("Status changed to up", "target", target)
		}
	}
}

// A change from three conscutive -1 to valid RTT means status change to up
func isStatusToUp(rttList []time.Duration) bool {
	if l := len(rttList); l > 4 {
		if rttList[l-2] == -1 && rttList[l-3] == -1 && rttList[l-4] == -1 {
			if rttList[l-1] != -1 {
				return true
			}
		}
	}
	return false
}
