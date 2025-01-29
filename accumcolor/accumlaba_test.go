// This file defines a suite of tests for accumcolor.AccumLabA.

package accumcolor

import (
	"image/color"
	"math"
	"testing"

	"github.com/lucasb-eyer/go-colorful"
)

// compareFloats aborts if two floats are not within some threshold distance.
func compareFloats(t *testing.T, nm string, act, exp float64) {
	const maxDiffAllowed = 1e-5
	if math.Abs(exp-act) > maxDiffAllowed {
		t.Fatalf("expected %s = %v but saw %v", nm, exp, act)
	}
}

// TestLabAValid ensures we can distinguish valid from invalid colors.
func TestLabAValid(t *testing.T) {
	var c AccumLabA
	if !c.Valid() {
		t.Fatalf("expected %v to be valid, but it is deemed invalid", c)
	}
	c.A = 0.5
	if c.Valid() {
		t.Fatalf("expected %v to be invalid, but it is deemed valid", c)
	}
	c.Tally = 1
	if !c.Valid() {
		t.Fatalf("expected %v to be valid, but it is deemed invalid", c)
	}
	c = AccumLabA{
		L:     1.0,
		A:     -1.0,
		B:     1.0,
		Alpha: 255,
		Tally: 1,
	}
	if !c.Valid() {
		t.Fatalf("expected %v to be valid, but it is deemed invalid", c)
	}
	c.B += 0.25
	if c.Valid() {
		t.Fatalf("expected %v to be invalid, but it is deemed valid", c)
	}
	c.B = 1.0
	c.A *= 2.0
	if c.Valid() {
		t.Fatalf("expected %v to be invalid, but it is deemed valid", c)
	}
	c.Tally = 2
	if !c.Valid() {
		t.Fatalf("expected %v to be valid, but it is deemed invalid", c)
	}
	const big = 34359738641 // Larger than 2^32
	c = AccumLabA{
		L:     1.0 * big,
		A:     1.0 * big,
		B:     -1.0 * big,
		Alpha: 255 * big,
		Tally: big,
	}
	if !c.Valid() {
		t.Fatalf("expected %v to be valid, but it is deemed invalid", c)
	}
}

// TestLabAAdd ensures adding a number of colors produces the expected total.
func TestLabAAdd(t *testing.T) {
	// Add up a number of colors.
	const n = 5
	var acc AccumLabA
	var exp float64
	for i := 0; i < n; i++ {
		v := float64(i) / float64(n*2)
		clr := colorful.Lab(v, v, -v)
		acc.Add(clr)
		exp += v
	}

	// Check that the totals are as expected.
	compareFloats(t, "L", acc.L, exp)
	compareFloats(t, "A", acc.A, exp)
	compareFloats(t, "B", acc.B, -exp)
	if acc.Alpha != n*255 {
		t.Fatalf("expected Alpha = %d but saw %d", n*255, acc.Alpha)
	}
	if acc.Tally != n {
		t.Fatalf("expected Tally = %d but saw %d", n, acc.Tally)
	}
}

// TestLabAAverage ensures that averaging colors produces the expected result.
func TestLabAAverage(t *testing.T) {
	// Test repeatedly adding a color to itself.
	exp := AccumLabA{
		L:     0.25,
		A:     -0.75,
		B:     0.50,
		Alpha: 250,
		Tally: 1,
	}
	var sum AccumLabA
	for i := 0; i < 10; i++ {
		sum.Add(exp)
	}
	act := sum.Average()
	compareFloats(t, "L", act.L, exp.L)
	compareFloats(t, "A", act.A, exp.A)
	compareFloats(t, "B", act.B, exp.B)
	if act.Alpha != exp.Alpha {
		t.Fatalf("expected Alpha = %d but saw %d", exp.Alpha, act.Alpha)
	}
	if act.Tally != 1 {
		t.Fatalf("expected Tally = 1 but saw %d", act.Tally)
	}
}

// TestLabAConvert ensures that we can convert to and from an AccumLabA.
func TestLabAConvert(t *testing.T) {
	rgba := color.RGBA{
		R: 0x22,
		G: 0x44,
		B: 0x66,
		A: 0x88,
	}
	laba := AccumLabAModel.Convert(rgba).(AccumLabA)
	rgba2 := color.RGBAModel.Convert(laba).(color.RGBA)
	if rgba != rgba2 {
		t.Fatalf("expected RGBA = %v but saw %v", rgba, rgba2)
	}
}

// TestLabAScale ensures that a weighted sum of colors produces the expected
// total.
func TestLabAScale(t *testing.T) {
	// Define pure red, green, and blue RGB colors.
	convertRGB := func(r, g, b uint8) AccumLabA {
		clr := color.RGBA{R: r, G: g, B: b, A: 255}
		return AccumLabAModel.Convert(clr).(AccumLabA)
	}
	red := convertRGB(255, 0, 0)
	green := convertRGB(0, 255, 0)
	blue := convertRGB(0, 0, 255)

	// Average (2*red + 1*green + 0*blue)/3.
	red.Scale(2)
	green.Scale(1)
	blue.Scale(0)
	var sum AccumLabA
	sum.Add(red)
	sum.Add(green)
	sum.Add(blue)
	avg := sum.Average()

	// Ensure we wound up with a proper orange.
	nrgba := color.NRGBAModel.Convert(avg).(color.NRGBA)
	orange := color.NRGBA{R: 223, G: 138, B: 0, A: 255}
	if nrgba != orange {
		t.Fatalf("expected %v but saw %v", orange, nrgba)
	}
}
