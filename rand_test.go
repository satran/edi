package main

import (
	"testing"
)

func TestRandString(t *testing.T) {
	r := RandString(6)
	t.Log(r)
	if len(r) != 6 {
		t.Errorf(r)
	}
}
