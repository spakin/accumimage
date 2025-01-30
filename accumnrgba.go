// This file defines the AccumNRGBA type and associated methods.

package accumimage

import (
	"image"
	"image/color"
	"math/bits"

	"github.com/spakin/accumimage/accumcolor"
)

// NOTE: Many of the functions and methods in this file were copied verbatim or
// nearly verbatim from the Go standard library (image/image.go).

// An AccumNRGBA is an in-memory image whose At method returns accumcolor.NRGBA
// values.
type AccumNRGBA struct {
	// Pix holds the image's pixels, in R, G, B, A, Tally order. The pixel
	// at (x, y) starts at Pix[(y-Rect.Min.Y)*Stride + (x-Rect.Min.X)*5].
	Pix []uint64
	// Stride is the Pix stride (in uint64s) between vertically adjacent
	// pixels.
	Stride int
	// Rect is the image's bounds.
	Rect image.Rectangle
}

// mul3NonNeg returns (x * y * z), unless at least one argument is negative or
// if the computation overflows the int type, in which case it returns -1.
func mul3NonNeg(x int, y int, z int) int {
	if (x < 0) || (y < 0) || (z < 0) {
		return -1
	}
	hi, lo := bits.Mul64(uint64(x), uint64(y))
	if hi != 0 {
		return -1
	}
	hi, lo = bits.Mul64(lo, uint64(z))
	if hi != 0 {
		return -1
	}
	a := int(lo)
	if (a < 0) || (uint64(a) != lo) {
		return -1
	}
	return a
}

// pixelBufferLength returns the length of the []uint8 typed Pix slice field
// for the NewXxx functions. Conceptually, this is just (bpp * width * height),
// but this function panics if at least one of those is negative or if the
// computation would overflow the int type.
func pixelBufferLength(bytesPerPixel int, r image.Rectangle, imageTypeName string) int {
	totalLength := mul3NonNeg(bytesPerPixel, r.Dx(), r.Dy())
	if totalLength < 0 {
		panic("accumimage: New" + imageTypeName + " Rectangle has huge or negative dimensions")
	}
	return totalLength
}

// NewAccumNRGBA returns a new AccumNRGBA image with the given bounds.
func NewAccumNRGBA(r image.Rectangle) *AccumNRGBA {
	return &AccumNRGBA{
		Pix:    make([]uint64, pixelBufferLength(5, r, "AccumNRGBA")),
		Stride: 5 * r.Dx(),
		Rect:   r,
	}
}

// At returns the color of the pixel at (x, y) as a color.Color.
func (p *AccumNRGBA) At(x, y int) color.Color {
	return p.AccumNRGBAAt(x, y)
}

// AccumNRGBAAt returns the color of the pixel at (x, y) as an
// accumcolor.AccumNRGBA.
func (p *AccumNRGBA) AccumNRGBAAt(x, y int) accumcolor.AccumNRGBA {
	if !(image.Point{x, y}.In(p.Rect)) {
		return accumcolor.AccumNRGBA{}
	}
	i := p.PixOffset(x, y)
	s := p.Pix[i : i+5 : i+5] // Small cap improves performance, see https://golang.org/issue/27857
	return accumcolor.AccumNRGBA{R: s[0], G: s[1], B: s[2], A: s[3], Tally: s[4]}
}

// NRGBAAt returns the color of the pixel at (x, y) as a color.NRGBA.
func (p *AccumNRGBA) NRGBAAt(x, y int) color.NRGBA {
	c := p.AccumNRGBAAt(x, y)
	return color.NRGBAModel.Convert(c).(color.NRGBA)
}

// PixOffset returns the index of the first element of Pix that corresponds to
// the pixel at (x, y).
func (p *AccumNRGBA) PixOffset(x, y int) int {
	return (y-p.Rect.Min.Y)*p.Stride + (x-p.Rect.Min.X)*5
}

// Bounds returns the domain for which At can return non-zero color.
func (p *AccumNRGBA) Bounds() image.Rectangle { return p.Rect }

// ColorModel returns the AccumNRGBA's color model (always
// accumcolor.AccumNRGBAModel).
func (p *AccumNRGBA) ColorModel() color.Model {
	return accumcolor.AccumNRGBAModel
}

// Opaque scans the entire image and reports whether it is fully opaque.
func (p *AccumNRGBA) Opaque() bool {
	if p.Rect.Empty() {
		return true
	}
	i0, i1 := 3, p.Rect.Dx()*5
	for y := p.Rect.Min.Y; y < p.Rect.Max.Y; y++ {
		for i := i0; i < i1; i += 5 {
			tally := p.Pix[i+1]
			if tally == 0 {
				return false // No color at this position
			}
			if p.Pix[i] != 0xff*tally {
				return false // Not fully opaque
			}
		}
		i0 += p.Stride
		i1 += p.Stride
	}
	return true
}

// RGBA64At returns the color of the pixel at (x, y) as a color.RGBA64.
func (p *AccumNRGBA) RGBA64At(x, y int) color.RGBA64 {
	r, g, b, a := p.AccumNRGBAAt(x, y).RGBA()
	return color.RGBA64{uint16(r), uint16(g), uint16(b), uint16(a)}
}

// Set sets the pixel at (x, y) to a given color of any type.
func (p *AccumNRGBA) Set(x, y int, c color.Color) {
	if !(image.Point{x, y}.In(p.Rect)) {
		return
	}
	i := p.PixOffset(x, y)
	c1 := accumcolor.AccumNRGBAModel.Convert(c).(accumcolor.AccumNRGBA)
	s := p.Pix[i : i+5 : i+5] // Small cap improves performance, see https://golang.org/issue/27857
	s[0] = c1.R
	s[1] = c1.G
	s[2] = c1.B
	s[3] = c1.A
	s[4] = c1.Tally
}

// Add accumulates a given color of any type to the pixel at (x, y).
func (p *AccumNRGBA) Add(x, y int, c color.Color) {
	if !(image.Point{x, y}.In(p.Rect)) {
		return
	}
	i := p.PixOffset(x, y)
	c1 := accumcolor.AccumNRGBAModel.Convert(c).(accumcolor.AccumNRGBA)
	s := p.Pix[i : i+5 : i+5] // Small cap improves performance, see https://golang.org/issue/27857
	s[0] += c1.R
	s[1] += c1.G
	s[2] += c1.B
	s[3] += c1.A
	s[4] += c1.Tally
}

// SetAccumNRGBA sets the pixel at (x, y) to a given color of type
// accumcolor.AccumNRGBA.
func (p *AccumNRGBA) SetAccumNRGBA(x, y int, c accumcolor.AccumNRGBA) {
	if !(image.Point{x, y}.In(p.Rect)) {
		return
	}
	i := p.PixOffset(x, y)
	s := p.Pix[i : i+5 : i+5] // Small cap improves performance, see https://golang.org/issue/27857
	s[0] = c.R
	s[1] = c.G
	s[2] = c.B
	s[3] = c.A
	s[4] = c.Tally
}

// AddAccumNRGBA accumulates a given color of type accumcolor.AccumNRGBA to the
// pixel at (x, y).
func (p *AccumNRGBA) AddAccumNRGBA(x, y int, c accumcolor.AccumNRGBA) {
	if !(image.Point{x, y}.In(p.Rect)) {
		return
	}
	i := p.PixOffset(x, y)
	s := p.Pix[i : i+5 : i+5] // Small cap improves performance, see https://golang.org/issue/27857
	s[0] += c.R
	s[1] += c.G
	s[2] += c.B
	s[3] += c.A
	s[4] += c.Tally
}

// SetRGBA64 sets the pixel at (x, y) to a given color of type color.RGBA64.
func (p *AccumNRGBA) SetRGBA64(x, y int, c color.RGBA64) {
	if !(image.Point{x, y}.In(p.Rect)) {
		return
	}
	r, g, b, a := uint32(c.R), uint32(c.G), uint32(c.B), uint32(c.A)
	if (a != 0) && (a != 0xffff) {
		r = (r * 0xffff) / a
		g = (g * 0xffff) / a
		b = (b * 0xffff) / a
	}
	i := p.PixOffset(x, y)
	s := p.Pix[i : i+5 : i+5] // Small cap improves performance, see https://golang.org/issue/27857
	s[0] = uint64(r >> 8)
	s[1] = uint64(g >> 8)
	s[2] = uint64(b >> 8)
	s[3] = uint64(a >> 8)
	s[4] = 1
}

// AddRGBA64 accumulates a given color of type color.RGBA64 to the pixel at (x, y).
func (p *AccumNRGBA) AddRGBA64(x, y int, c color.RGBA64) {
	if !(image.Point{x, y}.In(p.Rect)) {
		return
	}
	r, g, b, a := uint32(c.R), uint32(c.G), uint32(c.B), uint32(c.A)
	if (a != 0) && (a != 0xffff) {
		r = (r * 0xffff) / a
		g = (g * 0xffff) / a
		b = (b * 0xffff) / a
	}
	i := p.PixOffset(x, y)
	s := p.Pix[i : i+5 : i+5] // Small cap improves performance, see https://golang.org/issue/27857
	s[0] += uint64(r >> 8)
	s[1] += uint64(g >> 8)
	s[2] += uint64(b >> 8)
	s[3] += uint64(a >> 8)
	s[4]++
}

// SubImage returns an image representing the portion of the image p visible
// through r. The returned value shares pixels with the original image.
func (p *AccumNRGBA) SubImage(r image.Rectangle) image.Image {
	r = r.Intersect(p.Rect)
	// If r1 and r2 are Rectangles, r1.Intersect(r2) is not guaranteed to
	// be inside either r1 or r2 if the intersection is empty. Without
	// explicitly checking for this, the Pix[i:] expression below can
	// panic.
	if r.Empty() {
		return &AccumNRGBA{}
	}
	i := p.PixOffset(r.Min.X, r.Min.Y)
	return &AccumNRGBA{
		Pix:    p.Pix[i:],
		Stride: p.Stride,
		Rect:   r,
	}
}
