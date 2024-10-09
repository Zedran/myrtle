package main

import "testing"

// Tests whether the user input is parsed properly. No empty commands
// and duplicates are allowed.
func TestNewPhrase(t *testing.T) {
	case1 := NewPhrase("/c/c")

	if len(case1.Commands) != 1 || case1.Commands[0] != "c" {
		t.Fatalf("Case1: %s", case1)
	}

	case2 := NewPhrase("iss /")

	if len(case2.Commands) != 0 || case2.Object != "iss" {
		t.Fatalf("Case2: %s; %v", case2.Object, case2.Commands)
	}
}
