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
		slog.Debug("Target read", "target", target)
		if isIP(target) || isHostname(target) {
			targets = append(targets, target)
		} else {
			slog.Warn("Skipping invalid target", "target", target)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	slog.Debug("Target list completed", "targets", targets)
	return targets, nil
}
