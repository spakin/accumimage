// This file defines a suite of tests for accumcolor.AccumNRGBA.

package accumcolor

import (
	"image/color"
	"testing"
)

// TestAccumNRGBAAdd ensures adding a number of colors produces the expected
// total.
func TestAccumNRGBAAdd(t *testing.T) {
	// Add up a large number of colors.
	const n = 5
	var acc AccumNRGBA
	for r := uint8(0); r <= n-4; r++ {
		g := r + 1
		b := g + 1
		a := b + 1
		clr := color.NRGBA{
			R: r,
			G: g,
			B: b,
			A: a,
		}
		acc.Add(clr)
	}

	// Check that the totals are as expected.
	var exp uint64 = ((n - 4) * (n - 3)) / 2
	if acc.R != exp {
		t.Fatalf("expected R = %d but saw %d", exp, acc.R)
	}
	exp += (n - 3)
	if acc.G != exp {
		t.Fatalf("expected G = %d but saw %d", exp, acc.G)
	}
	exp += (n - 2) - 1
	if acc.B != exp {
		t.Fatalf("expected B = %d but saw %d", exp, acc.B)
	}
	exp += (n - 1) - 2
	if acc.A != exp {
		t.Fatalf("expected A = %d but saw %d", exp, acc.A)
	}
	exp = (n - 4) + 1
	if acc.Tally != exp {
		t.Fatalf("expected Tally = %d but saw %d", exp, acc.Tally)
	}
}
