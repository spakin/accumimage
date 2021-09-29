// This file defines a suite of tests for accumimage.AccumNRGBA.

package accumimage

import (
	"image"
	"image/color"
	"testing"
)

// TestAdd adds a number of colors together and checks the result.  It uses
// NewAccumNRGBA and AccumNRGBA's Add and AccumNRGBAAt methods.
func TestAdd(t *testing.T) {
	// Construct a row of colors.
	img := NewAccumNRGBA(image.Rect(0, 0, 256, 1))

	// Repeatedly add the same color to each pixel.
	const n = 10
	for j := 0; j < n; j++ {
		for i := 0; i < 256; i++ {
			c := color.NRGBA{
				R: uint8(i),
				G: uint8(i),
				B: uint8(i),
				A: uint8(i),
			}
			img.Add(i, 0, c)
		}
	}

	// Confirm that each value is as expected.
	for i := uint64(0); i < 256; i++ {
		c := img.AccumNRGBAAt(int(i), 0)
		if c.Tally != n {
			t.Fatalf("incorrect tally at position (%d, 0)", i)
		}
		if c.R != i*n || c.G != i*n || c.B != i*n || c.A != i*n {
			t.Fatalf("incorrect color at position (%d, 0)", i)
		}
	}
}

// TestNRGBA adds together different numbers of colors and checks that the
// averages are as expected.
func TestNRGBA(t *testing.T) {
	// Construct a column of colors.
	const n = 100
	img := NewAccumNRGBA(image.Rect(0, 0, 1, n))

	// Accumulate the most colors to the first pixel, less to the second,
	// less to the third, and so forth.
	for i := 0; i < n; i++ {
		c := color.NRGBA{
			R: uint8(i),
			G: uint8((i + 10) % 256),
			B: uint8((i + 20) % 256),
			A: uint8((i + 30) % 256),
		}
		for j := 0; j <= i; j++ {
			img.Add(0, j, c)
		}
	}

	// Confirm that each pixel contains the expected color.
	for i := 0; i < n; i++ {
		c := img.NRGBAAt(0, i)
		base := uint8((n + i - 1) / 2)
		exp := color.NRGBA{
			R: base,
			G: base + 10,
			B: base + 20,
			A: base + 30,
		}
		if c != exp {
			t.Fatalf("expected %v but saw %v", exp, c)
		}
	}
}
