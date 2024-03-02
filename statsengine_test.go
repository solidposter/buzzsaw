package main

import (
	"reflect"
	"testing"
	"time"
)

func TestStatusToDown(t *testing.T) {
	// Test cases
	tests := []struct {
		name    string
		targets map[string][]time.Duration
		want    []string
	}{
		{
			name: "Targets valid",
			targets: map[string][]time.Duration{
				"targetIsDown":        targetIsDown,
				"targetIsUp":          targetIsUp,
				"targetIsAlreadyDown": targetIsAlreadyDown,
				"targetIsAlreadyUp":   targetIsAlreadyUp,
				"targetToDown":        targetToDown,
				"targetToUp":          targetToUp,
				"targetFlappyToDown":  targetFlappyToDown,
				"targetFlappyToUp":    targetFlappyToUp,
			},
			want: []string{"targetToDown", "targetFlappyToDown"},
		},
		{
			name: "Targets too short",
			targets: map[string][]time.Duration{
				"empty":    targetEmpty,
				"toShoort": targetTooShort,
			},
			want: []string{},
		},
		// Add more test cases as needed
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := statusToDown(tc.targets)
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("statusToDown(%v) = %v; want %v", tc.targets, got, tc.want)
			}
		})
	}
}

func TestStatusToUp(t *testing.T) {
	// Test cases
	tests := []struct {
		name    string
		targets map[string][]time.Duration
		want    []string
	}{
		{
			name: "Targets valid",
			targets: map[string][]time.Duration{
				"targetIsDown":        targetIsDown,
				"targetIsUp":          targetIsUp,
				"targetIsAlreadyDown": targetIsAlreadyDown,
				"targetIsAlreadyUp":   targetIsAlreadyUp,
				"targetToDown":        targetToDown,
				"targetToUp":          targetToUp,
				"targetFlappyToDown":  targetFlappyToDown,
				"targetFlappyToUp":    targetFlappyToUp,
			},
			want: []string{"targetToUp", "targetFlappyToUp"},
		},
		{
			name: "Targets too short",
			targets: map[string][]time.Duration{
				"empty":    targetEmpty,
				"toShoort": targetTooShort,
			},
			want: []string{},
		},
		// Add more test cases as needed
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := statusToUp(tc.targets)
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("statusToDown(%v) = %v; want %v", tc.targets, got, tc.want)
			}
		})
	}
}
