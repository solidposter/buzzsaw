package main

import (
	"bufio"
	"os"
)

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
