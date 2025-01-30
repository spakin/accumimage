// This file defines the AccumLabA type and associated methods.

package accumimage

import (
	"image"
	"image/color"

	"github.com/lucasb-eyer/go-colorful"
	"github.com/spakin/accumimage/accumcolor"
)

// An AccumLabA is an in-memory image whose At method returns
// accumcolor.AccumLabA values.
type AccumLabA struct {
	// Pix holds the image's pixels.  Pix[0][0] corresponds to (Rect.Min.X,
	// Rect.Min.Y).
	Pix [][]accumcolor.AccumLabA
	// Rect is the image's bounds.
	Rect image.Rectangle
}

// NewAccumLabA returns a new AccumLabA image with the given bounds.
func NewAccumLabA(r image.Rectangle) *AccumLabA {
	wd, ht := r.Dx(), r.Dy()
	pixels := make([][]accumcolor.AccumLabA, ht)
	for r := range pixels {
		pixels[r] = make([]accumcolor.AccumLabA, wd)
	}
	return &AccumLabA{
		Pix:  pixels,
		Rect: r,
	}
}

// At returns the color of the pixel at (x, y) as a color.Color.
func (p *AccumLabA) At(x, y int) color.Color {
	return p.AccumLabAAt(x, y)
}

// AccumLabAAt returns the color of the pixel at (x, y) as an
// accumcolor.AccumLabA.
func (p *AccumLabA) AccumLabAAt(x, y int) accumcolor.AccumLabA {
	if !(image.Point{x, y}.In(p.Rect)) {
		return accumcolor.AccumLabA{}
	}
	return p.Pix[y-p.Rect.Min.Y][x-p.Rect.Min.X]
}

// ColorfulAt returns the color of the pixel at (x, y) as a fully opaque
// colorful.Color.
func (p *AccumLabA) ColorfulAt(x, y int) colorful.Color {
	c := p.AccumLabAAt(x, y).Average()
	return colorful.Lab(c.L, c.A, c.B)
}

// Bounds returns the domain for which At can return non-zero color.
func (p *AccumLabA) Bounds() image.Rectangle { return p.Rect }

// ColorModel returns the AccumLabA's color model (always
// accumcolor.AccumLabAModel).
func (p *AccumLabA) ColorModel() color.Model {
	return accumcolor.AccumLabAModel
}

// Opaque scans the entire image and reports whether it is fully opaque.
func (p *AccumLabA) Opaque() bool {
	if p.Rect.Empty() {
		return true
	}
	for _, row := range p.Pix {
		for _, clr := range row {
			if clr.Alpha != 255*clr.Tally {
				return false
			}
		}
	}
	return true
}

// RGBA64At returns the color of the pixel at (x, y) as a color.RGBA64.
func (p *AccumLabA) RGBA64At(x, y int) color.RGBA64 {
	r, g, b, a := p.AccumLabAAt(x, y).RGBA()
	return color.RGBA64{uint16(r), uint16(g), uint16(b), uint16(a)}
}

// Set sets the pixel at (x, y) to a given color of any type.
func (p *AccumLabA) Set(x, y int, c color.Color) {
	if !(image.Point{x, y}.In(p.Rect)) {
		return
	}
	clr := accumcolor.AccumLabAModel.Convert(c).(accumcolor.AccumLabA)
	p.Pix[y-p.Rect.Min.Y][x-p.Rect.Min.X] = clr
}

// Add accumulates a given color of any type to the pixel at (x, y).
func (p *AccumLabA) Add(x, y int, c color.Color) {
	if !(image.Point{x, y}.In(p.Rect)) {
		return
	}
	clr := accumcolor.AccumLabAModel.Convert(c).(accumcolor.AccumLabA)
	p.Pix[y-p.Rect.Min.Y][x-p.Rect.Min.X].Add(clr)
}

// SetAccumLabA sets the pixel at (x, y) to a given color of type
// accumcolor.AccumLabA.
func (p *AccumLabA) SetAccumLabA(x, y int, c accumcolor.AccumLabA) {
	if !(image.Point{x, y}.In(p.Rect)) {
		return
	}
	p.Pix[y-p.Rect.Min.Y][x-p.Rect.Min.X] = c
}

// AddAccumLabA accumulates a given color of type accumcolor.AccumLabA to the
// pixel at (x, y).
func (p *AccumLabA) AddAccumLabA(x, y int, c accumcolor.AccumLabA) {
	if !(image.Point{x, y}.In(p.Rect)) {
		return
	}
	p.Pix[y-p.Rect.Min.Y][x-p.Rect.Min.X].Add(c)
}

// SetRGBA64 sets the pixel at (x, y) to a given color of type color.RGBA64.
func (p *AccumLabA) SetRGBA64(x, y int, c color.RGBA64) {
	if !(image.Point{x, y}.In(p.Rect)) {
		return
	}
	if c.A == 0 {
		p.Pix[y-p.Rect.Min.Y][x-p.Rect.Min.X] = accumcolor.AccumLabA{Tally: 1}
		return
	}
	alpha := float64(c.A)
	clr := colorful.Color{
		R: float64(c.R) / alpha,
		G: float64(c.G) / alpha,
		B: float64(c.B) / alpha,
	}
	L, a, b := clr.Lab()
	p.Pix[y-p.Rect.Min.Y][x-p.Rect.Min.X] = accumcolor.AccumLabA{
		L:     L,
		A:     a,
		B:     b,
		Alpha: uint64(c.A >> 8),
		Tally: 1,
	}
}

// AddRGBA64 accumulates a given color of type color.RGBA64 to the pixel at (x, y).
func (p *AccumLabA) AddRGBA64(x, y int, c color.RGBA64) {
	if !(image.Point{x, y}.In(p.Rect)) {
		return
	}
	if c.A == 0 {
		p.Pix[y-p.Rect.Min.Y][x-p.Rect.Min.X].Add(accumcolor.AccumLabA{Tally: 1})
		return
	}
	alpha := float64(c.A)
	clr := colorful.Color{
		R: float64(c.R) / alpha,
		G: float64(c.G) / alpha,
		B: float64(c.B) / alpha,
	}
	L, a, b := clr.Lab()
	p.Pix[y-p.Rect.Min.Y][x-p.Rect.Min.X].Add(accumcolor.AccumLabA{
		L:     L,
		A:     a,
		B:     b,
		Alpha: uint64(c.A >> 8),
		Tally: 1,
	})
}

// SubImage returns an image representing the portion of the image p visible
// through r. The returned value shares pixels with the original image.
func (p *AccumLabA) SubImage(r image.Rectangle) image.Image {
	r = r.Intersect(p.Rect)
	// If r1 and r2 are Rectangles, r1.Intersect(r2) is not guaranteed to
	// be inside either r1 or r2.  We explicitly check for this situation.
	if r.Empty() {
		return &AccumLabA{}
	}

	// Create new rows, each containing a subset of the original columns.
	wd, ht := r.Dx(), r.Dy()
	xOfs := r.Min.X - p.Rect.Min.X
	yOfs := r.Min.Y - p.Rect.Min.Y
	pixels := make([][]accumcolor.AccumLabA, ht)
	for r := range pixels {
		pixels[r] = p.Pix[yOfs+r][xOfs : xOfs+wd]
	}
	return &AccumLabA{
		Pix:  pixels,
		Rect: r,
	}
}
