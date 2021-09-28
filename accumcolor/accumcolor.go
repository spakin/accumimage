/*
accumcolor provides support for colors that can be accumulated and averaged.

An AccumNRGBA is analogous to a color.NRGBA but additionally supports adding
colors together.  It maintains a tally of the total number of colors that have
been accumulated.  The Add method accumulates more color onto an existing
AccumNRGBA.  The NRGBA method returns the average color of the entire
accumulation as a color.NRGBa.
*/
package accumcolor

import "image/color"

// An AccumNRGBA is a color.Color that supports accumulation of
// non-alpha-premultiplied RGBA color values.
type AccumNRGBA struct {
	R     uint64
	G     uint64
	B     uint64
	A     uint64
	Tally uint64
}

// RGBA convers an AccumNRGBA to alpha-premultiplied colors.
func (c AccumNRGBA) RGBA() (r, g, b, a uint32) {
	if c.Tally == 0 {
		return
	}
	nrgba := color.NRGBA{
		R: uint8(c.R / c.Tally),
		G: uint8(c.G / c.Tally),
		B: uint8(c.B / c.Tally),
		A: uint8(c.A / c.Tally),
	}
	r, g, b, a = nrgba.RGBA()
	return
}

// accumNRGBAModel is used to define a color model for AccumNRGBA.
func accumNRGBAModel(c color.Color) color.Color {
	if _, ok := c.(AccumNRGBA); ok {
		return c
	}
	nrgba := color.NRGBAModel.Convert(c).(color.NRGBA)
	return AccumNRGBA{
		R:     uint64(nrgba.R),
		G:     uint64(nrgba.G),
		B:     uint64(nrgba.B),
		A:     uint64(nrgba.A),
		Tally: 1,
	}
}

// AccumNRGBAModel converts any color.Color to an AccumNRGBA color.
var AccumNRGBAModel = color.ModelFunc(accumNRGBAModel)

// Add accumulates color.
func (c *AccumNRGBA) Add(clr color.Color) {
	other := AccumNRGBAModel.Convert(clr).(AccumNRGBA)
	c.R += other.R
	c.G += other.G
	c.B += other.B
	c.A += other.A
	c.Tally += other.Tally
}

// NRGBA averages the accumulated color of an AccumNRGBA to produce an ordinary
// color.NRGBA.
func (c AccumNRGBA) NRGBA() color.NRGBA {
	if c.Tally == 0 {
		return color.NRGBA{}
	}
	return color.NRGBA{
		R: uint8(c.R / c.Tally),
		G: uint8(c.G / c.Tally),
		B: uint8(c.B / c.Tally),
		A: uint8(c.A / c.Tally),
	}
}
