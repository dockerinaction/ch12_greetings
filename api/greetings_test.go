package main

import (
	"math/rand"
	"testing"
)

func TestSelectRandom(t *testing.T) {
	someStrings := []string{
		"abcd",
		"efgh",
		"jklmn",
		"opqrs",
		"tuvwx",
		"yz",
	}

	expectedStrings := []string{
		"yz",
		"yz",
		"jklmn",
		"abcd",
		"efgh",
		"efgh",
		"opqrs",
		"jklmn",
		"jklmn",
		"efgh",
	}

	r = rand.New(rand.NewSource(42))

	for _, expected := range expectedStrings {
		t.Run(expected, func(t *testing.T) {
			actual := SelectRandom(someStrings, r)
			if actual != expected {
				t.Errorf("got %q, want %q", actual, expected)
			}
		})
	}
}
