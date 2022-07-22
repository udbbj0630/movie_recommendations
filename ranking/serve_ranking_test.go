package main

import "testing"

func TestAdd(t *testing.T) {
	got := Add(1,4)
	want := 3

	if got != want 	{
		t.Fatalf(" expected %d, but got %d", want, got)
	}
}