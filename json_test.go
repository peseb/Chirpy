package main

import "testing"

func TestSimpleInput(t *testing.T) {
	input := "This is a kerfuffle opinion I need to share with the world"
	cleaned_input := cleanInput(input)

	expected := "This is a **** opinion I need to share with the world"
	if cleaned_input != expected {
		t.Errorf("cleaned_input = %s; exptected %s", cleaned_input, expected)
	}
}
