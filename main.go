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
	"flag"
	"fmt"
	"log/slog"
	"os"
	"time"
)

var version string // populated at build time

func main() {
	debugPtr := flag.Bool("d", false, "enable debug logging")
	filePtr := flag.String("f", "targets.txt", "targets file")
	logPtr := flag.String("l", "buzzsaw.log", "logfile")
	versPtr := flag.Bool("V", false, "print version info")
	flag.Parse()

	if *versPtr {
		fmt.Println("Version:", version)
		os.Exit(0)
	}

	slogsetup(*logPtr, *debugPtr)
	slog.Info("Starting", "version", version)

	targets, err := readTargets(*filePtr)
	if err != nil {
		slog.Error("Error reading targets", "error", err)
		os.Exit(1)
	}

	dispatcher := newDispatcher(len(targets))
	go dispatcher.start()
	statsEngine := newStatsEngine(len(targets))
	go statsEngine.start()
	go listener(dispatcher.getInput())

	// create and prep the pingers
	var pingers []*pinger
	for _, v := range targets {
		c, err := newPinger(v)
		if err != nil {
			slog.Warn("Skipping invalid target", "error", err)
			continue
		}
		if err := dispatcher.addPinger(c.getTarget(), c.getInput()); err != nil {
			slog.Warn("Skipping duplicate target", "error", err)
			continue
		}
		pingers = append(pingers, c)
	}
	// start the pingers
	for _, v := range pingers {
		time.Sleep(10 * time.Millisecond)
		go v.start(statsEngine.getInput())
	}
	slog.Info("Pingers started")

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
