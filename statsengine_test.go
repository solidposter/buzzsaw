package main

import (
	"reflect"
	"testing"
	"time"
)

func TestStatusToDown(t *testing.T) {
	tests := []struct {
		name    string
		targets map[string][]time.Duration
		want    []string
	}{
		{
			name: "Match targetToDown",
			targets: map[string][]time.Duration{
				"targetIsDown":        targetIsDown,
				"targetIsUp":          targetIsUp,
				"targetIsAlreadyDown": targetIsAlreadyDown,
				"targetIsAlreadyUp":   targetIsAlreadyUp,
				"targetToDown":        targetToDown,
				"targetToUp":          targetToUp,
				//"targetFlappyToDown":  targetFlappyToDown,
				"targetFlappyToUp": targetFlappyToUp,
				"empty":            targetEmpty,
				"toShoort":         targetTooShort,
			},
			want: []string{"targetToDown"},
		},
		{
			name: "Match targetFlappyToDown",
			targets: map[string][]time.Duration{
				"targetIsDown":        targetIsDown,
				"targetIsUp":          targetIsUp,
				"targetIsAlreadyDown": targetIsAlreadyDown,
				"targetIsAlreadyUp":   targetIsAlreadyUp,
				//"targetToDown":        targetToDown,
				"targetToUp":         targetToUp,
				"targetFlappyToDown": targetFlappyToDown,
				"targetFlappyToUp":   targetFlappyToUp,
				"empty":              targetEmpty,
				"toShoort":           targetTooShort,
			},
			want: []string{"targetFlappyToDown"},
		},
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
	tests := []struct {
		name    string
		targets map[string][]time.Duration
		want    []string
	}{
		{
			name: "Match targetToUp",
			targets: map[string][]time.Duration{
				"targetIsDown":        targetIsDown,
				"targetIsUp":          targetIsUp,
				"targetIsAlreadyDown": targetIsAlreadyDown,
				"targetIsAlreadyUp":   targetIsAlreadyUp,
				"targetToDown":        targetToDown,
				"targetToUp":          targetToUp,
				"targetFlappyToDown":  targetFlappyToDown,
				//"targetFlappyToUp": targetFlappyToUp,
				"empty":    targetEmpty,
				"toShoort": targetTooShort,
			},
			want: []string{"targetToUp"},
		},
		{
			name: "Match targetFlappyToUp",
			targets: map[string][]time.Duration{
				"targetIsDown":        targetIsDown,
				"targetIsUp":          targetIsUp,
				"targetIsAlreadyDown": targetIsAlreadyDown,
				"targetIsAlreadyUp":   targetIsAlreadyUp,
				"targetToDown":        targetToDown,
				//"targetToUp":         targetToUp,
				"targetFlappyToDown": targetFlappyToDown,
				"targetFlappyToUp":   targetFlappyToUp,
				"empty":              targetEmpty,
				"toShoort":           targetTooShort,
			},
			want: []string{"targetFlappyToUp"},
		},
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
