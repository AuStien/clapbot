package main

import (
	"testing"
)

func TestClap(t *testing.T) {
	testText := "this is just a simple test"
	clapped := addClap(testText)

	if clapped != "this :clap: is :clap: just :clap: a :clap: simple :clap: test :clap:" {
		t.Fatalf("expected %q, got %q", testText, clapped)
	}
}
