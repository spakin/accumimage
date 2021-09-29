// This file defines a suite of tests for accumimage.AccumNRGBA.

package accumimage

import (
	"image"
	"image/color"
	"testing"
)

// TestAdd adds a number of colors together and checks the result.
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
