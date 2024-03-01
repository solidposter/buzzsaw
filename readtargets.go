package main

import (
	"bufio"
	"log/slog"
	"os"
)

func readTargets(targetsfile string) ([]string, error) {
	targets := []string{}

	file, err := os.Open(targetsfile)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	slog.Debug("Targets file opened", "file", targetsfile)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		target := scanner.Text()
		if isIP(target) || isHostname(target) {
			targets = append(targets, target)
		} else {
			slog.Warn("Invalid target", "target", target)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	slog.Debug("Targets read", "targets", targets)
	return targets, nil
}
