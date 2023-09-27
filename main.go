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
	"bufio"
	"log"
	"os"
	"time"
)

func main() {

	targets, err := readTargets("targets.txt")
	if err != nil {
		log.Fatal(err)
	}

	dispatcher := newDispatcher()
	go dispatcher.start()
	statsEngine := newStatsEngine()
	go statsEngine.start()
	go startListener(dispatcher.getInput())

	// create and prep the pingers
	var pingers []*pinger
	for _, v := range targets {
		c := newPinger(v)
		dispatcher.addPinger(c.getTarget(), c.getInput())
		pingers = append(pingers, c)
	}
	// start the pingers
	for _, v := range pingers {
		time.Sleep(10 * time.Millisecond)
		go v.start(statsEngine.getInput())
	}

	// channel length monitor
	cm := newChanmonitor()
	cm.addDispatcher(dispatcher.getInput())   // listner -> dispatcher
	cm.addStatsengine(statsEngine.getInput()) // pingers -> statsengine
	for _, v := range pingers {               // dispatcher -> pingers
		cm.addPinger(v.getTarget(), v.getInput()) // one per pinger
	}
	go cm.start()

	<-(chan int)(nil) // wait forever
}

func readTargets(targetsfile string) ([]string, error) {
	targets := []string{}

	file, err := os.Open(targetsfile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		targets = append(targets, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return targets, nil
}
