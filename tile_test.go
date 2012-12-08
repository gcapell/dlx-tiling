package main

import (
	"fmt"
	"testing"
)

func TestASCIIToTile(test *testing.T) {
	fmt.Println("we're testing")

	src := "\n xx\nxx\n x\n x"
	expectedSquares := []square{
		{1, 0}, {2, 0},
		{0, 1}, {1, 1},
		{1, 2},
		{1, 3},
	}
	t, err := asciiToTile(src)
	if err != nil {
		test.Errorf("unexpected error %s with %q", err, src)
		return
	}
	if diff := diffSquares(expectedSquares, t.squares); diff != "" {
		test.Errorf("converting %q got %s", src, diff)
	}
}

func diffSquares(a, b []square) string {
	bag := make(map[square]bool)
	for _, s := range a {
		bag[s] = true
	}
	b_only := make([]square, 0)
	for _, s := range b {
		if bag[s] {
			delete(bag, s)
		} else {
			b_only = append(b_only, s)
		}
	}
	a_only := make([]square, 0)
	for k, _ := range bag {
		a_only = append(a_only, k)
	}
	problem := ""
	if len(a_only) != 0 {
		problem += fmt.Sprintf("only in a: %v", a_only)
	}
	if len(b_only) != 0 {
		problem += fmt.Sprintf("only in b: %v", b_only)
	}
	return problem
}
