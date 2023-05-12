package main

import (
	"fmt"
	"testing"
)

func TestClap(t *testing.T) {
	testText := "this is just a simple test"
	clapped := addClap(testText)

	if clapped != "this :clap: is :clap: just :clap: a :clap: simple :clap: test :clap:" {
		t.Fatalf("expected %q, got %q", testText, clapped)
	}
}

func TestRandomCase(t *testing.T) {
	seed := int64(1)
	testText := "this is just a simple test"
	randomCase := randomCase(testText, seed)

	fmt.Println(randomCase)

	if randomCase != "this Is JUSt a SimpLe TEst" {
		t.Fatalf("expected %q, got %q", testText, randomCase)
	}
}
