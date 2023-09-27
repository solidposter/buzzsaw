package main

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
