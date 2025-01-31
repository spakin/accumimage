// This file defines a suite of tests for accumcolor.NRGBA.

package accumcolor

import (
	"image/color"
	"testing"
)

// TestNRGBAValid ensures we can distinguish valid from invalid colors.
func TestNRGBAValid(t *testing.T) {
	var c NRGBA
	if !c.Valid() {
		t.Fatalf("expected %v to be valid, but it is deemed invalid", c)
	}
	c.G = 123
	if c.Valid() {
		t.Fatalf("expected %v to be invalid, but it is deemed valid", c)
	}
	c.Tally = 1
	if !c.Valid() {
		t.Fatalf("expected %v to be valid, but it is deemed invalid", c)
	}
	c = NRGBA{
		R:     255,
		G:     255,
		B:     255,
		A:     255,
		Tally: 1,
	}
	c.Tally = 1
	if !c.Valid() {
		t.Fatalf("expected %v to be valid, but it is deemed invalid", c)
	}
	c.B++
	if c.Valid() {
		t.Fatalf("expected %v to be invalid, but it is deemed valid", c)
	}
	c.R *= 2
	if c.Valid() {
		t.Fatalf("expected %v to be invalid, but it is deemed valid", c)
	}
	c.Tally = 2
	if !c.Valid() {
		t.Fatalf("expected %v to be valid, but it is deemed invalid", c)
	}
	const big = 34359738641 // Larger than 2^32
	c = NRGBA{
		R:     255 * big,
		G:     255 * big,
		B:     255 * big,
		A:     255 * big,
		Tally: big,
	}
	if !c.Valid() {
		t.Fatalf("expected %v to be valid, but it is deemed invalid", c)
	}
}

// TestNRGBAAdd ensures adding a number of colors produces the expected total.
func TestNRGBAAdd(t *testing.T) {
	// Add up a large number of colors.
	const n = 5
	var acc NRGBA
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

// TestNRGBA1 ensures that averaging colors produces the expected result.
func TestNRGBA1(t *testing.T) {
	// Test exclusively even numbers.
	c1 := NRGBA{
		R:     100,
		G:     110,
		B:     120,
		A:     130,
		Tally: 1,
	}
	c2 := NRGBA{
		R:     200,
		G:     210,
		B:     220,
		A:     230,
		Tally: 1,
	}
	var sumA NRGBA
	sumA.Add(c1)
	sumA.Add(c2)
	nrgba := sumA.NRGBA()
	exp := color.NRGBA{
		R: 150,
		G: 160,
		B: 170,
		A: 180,
	}
	if nrgba != exp {
		t.Fatalf("expected %v but saw %v", exp, nrgba)
	}

	// Confirm that halving an odd number rounds it down.
	c2 = NRGBA{
		R:     201,
		G:     211,
		B:     221,
		A:     231,
		Tally: 1,
	}
	var sumB NRGBA
	sumB.Add(c1)
	sumB.Add(c2)
	nrgba = sumB.NRGBA()
	if nrgba != exp {
		t.Fatalf("expected %v but saw %v", exp, nrgba)
	}
}

// TestNRGBA2 ensures that averaging colors produces the expected result.
func TestNRGBA2(t *testing.T) {
	c := color.NRGBA{
		R: 0,
		G: 128,
		B: 254,
		A: 255,
	}
	var sum NRGBA
	const n = 100000
	for i := 0; i < n; i++ {
		sum.Add(c)
	}
	nrgba := sum.NRGBA()
	if nrgba != c {
		t.Fatalf("expected %v but saw %v", c, nrgba)
	}
}

// TestNRGBARGBA ensures that we can convert an NRGBA to RGBA and back.
func TestNRGBARGBA(t *testing.T) {
	acc1 := NRGBA{
		R:     99,
		G:     100,
		B:     101,
		A:     255,
		Tally: 1,
	}
	r, g, b, a := acc1.RGBA()
	rgba := color.RGBA{
		R: uint8(r >> 8),
		G: uint8(g >> 8),
		B: uint8(b >> 8),
		A: uint8(a >> 8),
	}
	acc2 := NRGBAModel.Convert(rgba).(NRGBA)
	if acc2 != acc1 {
		t.Fatalf("expected %v but saw %v", acc1, acc2)
	}
}

// TestNRGBAScale ensures that a weighted sum of colors produces the expected
// total.
func TestNRGBAScale(t *testing.T) {
	// Define pure red, green, and blue colors.
	red := NRGBA{
		R:     255,
		G:     0,
		B:     0,
		A:     255,
		Tally: 1,
	}
	green := NRGBA{
		R:     0,
		G:     255,
		B:     0,
		A:     255,
		Tally: 1,
	}
	blue := NRGBA{
		R:     0,
		G:     0,
		B:     255,
		A:     255,
		Tally: 1,
	}

	// Average (2*red + 1*green + 0*blue)/3.
	red.Scale(2)
	green.Scale(1)
	blue.Scale(0)
	var sum NRGBA
	sum.Add(red)
	sum.Add(green)
	sum.Add(blue)

	// Ensure we wound up with dark orange.
	darkOrange := color.NRGBA{
		R: 170,
		G: 85,
		B: 0,
		A: 255,
	}
	nrgba := sum.NRGBA()
	if nrgba != darkOrange {
		t.Fatalf("expected %v but saw %v", darkOrange, nrgba)
	}
}
