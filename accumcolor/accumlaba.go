// This file defines the AccumLabA type and associated methods.

package accumcolor

import (
	"image/color"

	"github.com/lucasb-eyer/go-colorful"
)

// An AccumLabA is a color.Color that supports accumulation of CIE L*a*b* color
// values.  An invariant maintained by all methods is that either all fields
// are zero or each of L, A, B, and Alpha divided by Tally produces a value in
// its target range.
type AccumLabA struct {
	L     float64 // [0, 1]*Tally
	A     float64 // [-1, 1]*Tally
	B     float64 // [-1, 1]*Tally
	Alpha uint64  // [0, 255]*Tally
	Tally uint64
}

// Valid returns true if and only if an AccumLabA is valid.
func (c AccumLabA) Valid() bool {
	// The only time a Tally is allowed to be zero is if all other fields
	// are zero.
	if c.Tally == 0 {
		var zero AccumLabA
		return c == zero
	}

	// If Tally is nonzero, each other field divided by it must lie within
	// its target range.
	tally := float64(c.Tally)
	L := c.L / tally
	a := c.A / tally
	b := c.B / tally
	alpha := c.Alpha / c.Tally
	switch {

	case L < 0.0 || L > 1.0:
		return false
	case a < -1.0 || a > 1.0:
		return false
	case b < -1.0 || b > 1.0:
		return false
	case alpha > 255:
		return false
	default:
		return true
	}
}

// RGBA converts an AccumLabA to alpha-premultiplied colors.
func (c AccumLabA) RGBA() (r, g, b, a uint32) {
	if c.Tally == 0 {
		return
	}
	tally := float64(c.Tally)
	clr := colorful.Lab(c.L/tally, c.A/tally, c.B/tally).Clamped()
	alpha := float64(c.Alpha) / tally / 255.0
	r = uint32(clr.R*alpha*65535.0 + 0.5)
	g = uint32(clr.G*alpha*65535.0 + 0.5)
	b = uint32(clr.B*alpha*65535.0 + 0.5)
	a = uint32(alpha*65535.0 + 0.5)
	return
}

// accumLabAModel is used to define a color model for AccumLabA.
func accumLabAModel(c color.Color) color.Color {
	if _, ok := c.(AccumLabA); ok {
		return c
	}
	clr, _ := colorful.MakeColor(c)
	L, a, b := clr.Lab()
	_, _, _, alpha := c.RGBA()
	return AccumLabA{
		L:     L,
		A:     a,
		B:     b,
		Alpha: uint64(alpha >> 8),
		Tally: 1,
	}
}

// AccumLabAModel converts any color.Color to an AccumLabA color.
var AccumLabAModel = color.ModelFunc(accumLabAModel)

// Add accumulates color.
func (c *AccumLabA) Add(clr color.Color) {
	other := AccumLabAModel.Convert(clr).(AccumLabA)
	c.L += other.L
	c.A += other.A
	c.B += other.B
	c.Alpha += other.Alpha
	c.Tally += other.Tally
}

// Scale multiplies all components of an AccumLabA by a given value.  This
// does not change the effective color but can be used for performing weighted
// averages.
func (c *AccumLabA) Scale(w uint64) {
	w64 := float64(w)
	c.L *= w64
	c.A *= w64
	c.B *= w64
	c.Alpha *= w
	c.Tally *= w
}

// Average averages the accumulated color of an AccumLabA to produce an
// AccumLabA with a Tally of 1.
func (c AccumLabA) Average() AccumLabA {
	if c.Tally == 0 {
		return AccumLabA{}
	}
	tally := float64(c.Tally)
	return AccumLabA{
		L:     c.L / tally,
		A:     c.A / tally,
		B:     c.B / tally,
		Alpha: c.Alpha / c.Tally,
		Tally: 1,
	}
}

// Colorful averages the accumulated color of an AccumLabA to produce a
// colorful.Color (from the go-colorful package).
func (c AccumLabA) Colorful() colorful.Color {
	avg := c.Average()
	return colorful.Lab(avg.L, avg.A, avg.B)
}
