// This file defines a suite of tests for accumimage.LabA.

package accumimage

import (
	"image"
	"math"
	"testing"

	"github.com/spakin/accumimage/accumcolor"
)

// TestLabAAdd adds a number of colors together and checks the result.  It uses
// NewLabA and LabA's Add and LabAAt methods.
func TestLabAAdd(t *testing.T) {
	// Construct a row of colors.
	const wd = 256
	img := NewLabA(image.Rect(0, 0, wd, 1))

	// Repeatedly add the same color to each pixel.
	const n = 10
	c := accumcolor.LabA{
		L:     0.1,
		A:     0.2,
		B:     0.3,
		Alpha: 200,
		Tally: 1,
	}
	for j := 0; j < n; j++ {
		for i := 0; i < wd; i++ {
			img.Add(i, 0, c)
		}
	}

	// Confirm that each value is as expected.
	approxEqual := func(a, b float64) bool {
		return math.Abs(a-b) < 1e-6
	}
	for i := 0; i < wd; i++ {
		// Confirm the tally.
		c := img.LabAAt(int(i), 0)
		if c.Tally != n {
			t.Fatalf("incorrect tally at position (%d, 0)", i)
		}
		if c.Alpha != n*200 ||
			!approxEqual(c.L, 0.1*n) ||
			!approxEqual(c.A, 0.2*n) ||
			!approxEqual(c.B, 0.3*n) {
			t.Fatalf("incorrect color at position (%d, 0)", i)
		}
	}
}

// TestLabASubImage modifies values in a subimage and ensures these
// modifications are reflected in the original image.
func TestLabASubImage(t *testing.T) {
	// Construct a square image of a single color.
	const edge1 = 10
	const ofs1 = 3
	c1 := accumcolor.LabA{
		L:     0.8,
		A:     -0.5,
		B:     0.6,
		Alpha: 255,
		Tally: 1,
	}
	rect1 := image.Rect(ofs1, ofs1, edge1+ofs1, edge1+ofs1)
	img1 := NewLabA(rect1)
	for y := rect1.Min.Y; y < rect1.Max.Y; y++ {
		for x := rect1.Min.X; x < rect1.Max.X; x++ {
			img1.SetLabA(x, y, c1)
		}
	}

	// Assign all pixels in a square subimage a second color.
	const edge2 = 100 // Note: Much larger than edge1
	const ofs2 = 5
	c2 := accumcolor.LabA{
		L:     0.2,
		A:     0.7,
		B:     0.1,
		Alpha: 150,
		Tally: 1,
	}
	rect2 := image.Rect(ofs2, ofs2, edge2+ofs2, edge2+ofs2)
	img2 := img1.SubImage(rect2).(*LabA)
	rect2 = img2.Rect
	for y := rect2.Min.Y; y < rect2.Max.Y; y++ {
		for x := rect2.Min.X; x < rect2.Max.X; x++ {
			img2.SetLabA(x, y, c2)
		}
	}

	// Confirm that all pixels have the correct value.
	for y := rect1.Min.Y; y < rect1.Max.Y; y++ {
		for x := rect1.Min.X; x < rect1.Max.X; x++ {
			clr := img1.LabAAt(x, y)
			switch clr {
			case c1:
				if y >= ofs2 && x >= ofs2 {
					t.Fatalf("incorrect color C1 at (%d, %d)", x, y)
				}
			case c2:
				if y < ofs2 || x < ofs2 {
					t.Fatalf("incorrect color C2 at (%d, %d)", x, y)
				}
			default:
				t.Fatalf("unexpected color %v at (%d, %d)", clr, x, y)
			}
		}
	}
}
