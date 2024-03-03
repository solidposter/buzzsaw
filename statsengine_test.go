package main

import (
	"testing"
	"time"
)

func TestIsStatusToDown(t *testing.T) {
	var tests = []struct {
		expected bool
		rttList  []time.Duration
	}{
		{false, targetIsDown},
		{false, targetIsUp},
		{false, targetIsAlreadyDown},
		{false, targetIsAlreadyUp},
		{true, targetToDown},
		{false, targetToUp},
		{true, targetFlappyToDown},
		{false, targetFlappyToUp},
		{false, targetEmpty},
		{false, targetTooShort},
	}

	for _, test := range tests {
		if output := isStatusToDown(test.rttList); output != test.expected {
			t.Errorf("Output %t not equal to expected %t", output, test.expected)
		}
	}
}

func TestIsStatusToUp(t *testing.T) {
	var tests = []struct {
		expected bool
		rttList  []time.Duration
	}{
		{false, targetIsDown},
		{false, targetIsUp},
		{false, targetIsAlreadyDown},
		{false, targetIsAlreadyUp},
		{false, targetToDown},
		{true, targetToUp},
		{false, targetFlappyToDown},
		{true, targetFlappyToUp},
		{false, targetEmpty},
		{false, targetTooShort},
	}

	for _, test := range tests {
		if output := isStatusToUp(test.rttList); output != test.expected {
			t.Errorf("Output %t not equal to expected %t", output, test.expected)
		}
	}
}
