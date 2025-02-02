// This file defines the LabA type and associated methods.

package accumimage

import (
	"image"
	"image/color"

	"github.com/lucasb-eyer/go-colorful"
	"github.com/spakin/accumimage/v2/accumcolor"
)

// A LabA is an in-memory image whose At method returns accumcolor.LabA values.
type LabA struct {
	// Pix holds the image's pixels.  Pix[0][0] corresponds to (Rect.Min.X,
	// Rect.Min.Y).
	Pix [][]accumcolor.LabA
	// Rect is the image's bounds.
	Rect image.Rectangle
}

// NewLabA returns a new LabA image with the given bounds.
func NewLabA(r image.Rectangle) *LabA {
	wd, ht := r.Dx(), r.Dy()
	pixels := make([][]accumcolor.LabA, ht)
	for r := range pixels {
		pixels[r] = make([]accumcolor.LabA, wd)
	}
	return &LabA{
		Pix:  pixels,
		Rect: r,
	}
}

// At returns the color of the pixel at (x, y) as a color.Color.
func (p *LabA) At(x, y int) color.Color {
	return p.LabAAt(x, y)
}

// LabAAt returns the color of the pixel at (x, y) as an
// accumcolor.LabA.
func (p *LabA) LabAAt(x, y int) accumcolor.LabA {
	if !(image.Point{x, y}.In(p.Rect)) {
		return accumcolor.LabA{}
	}
	return p.Pix[y-p.Rect.Min.Y][x-p.Rect.Min.X]
}

// ColorfulAt returns the color of the pixel at (x, y) as a fully opaque
// colorful.Color (from the go-colorful package).
func (p *LabA) ColorfulAt(x, y int) colorful.Color {
	return p.LabAAt(x, y).Colorful()
}

// Bounds returns the domain for which At can return non-zero color.
func (p *LabA) Bounds() image.Rectangle { return p.Rect }

// ColorModel returns the LabA's color model (always
// accumcolor.LabAModel).
func (p *LabA) ColorModel() color.Model {
	return accumcolor.LabAModel
}

// Opaque scans the entire image and reports whether it is fully opaque.
func (p *LabA) Opaque() bool {
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
func (p *LabA) RGBA64At(x, y int) color.RGBA64 {
	r, g, b, a := p.LabAAt(x, y).RGBA()
	return color.RGBA64{uint16(r), uint16(g), uint16(b), uint16(a)}
}

// Set sets the pixel at (x, y) to a given color of any type.
func (p *LabA) Set(x, y int, c color.Color) {
	if !(image.Point{x, y}.In(p.Rect)) {
		return
	}
	clr := accumcolor.LabAModel.Convert(c).(accumcolor.LabA)
	p.Pix[y-p.Rect.Min.Y][x-p.Rect.Min.X] = clr
}

// Add accumulates a given color of any type to the pixel at (x, y).
func (p *LabA) Add(x, y int, c color.Color) {
	if !(image.Point{x, y}.In(p.Rect)) {
		return
	}
	clr := accumcolor.LabAModel.Convert(c).(accumcolor.LabA)
	p.Pix[y-p.Rect.Min.Y][x-p.Rect.Min.X].Add(clr)
}

// SetLabA sets the pixel at (x, y) to a given color of type
// accumcolor.LabA.
func (p *LabA) SetLabA(x, y int, c accumcolor.LabA) {
	if !(image.Point{x, y}.In(p.Rect)) {
		return
	}
	p.Pix[y-p.Rect.Min.Y][x-p.Rect.Min.X] = c
}

// AddLabA accumulates a given color of type accumcolor.LabA to the
// pixel at (x, y).
func (p *LabA) AddLabA(x, y int, c accumcolor.LabA) {
	if !(image.Point{x, y}.In(p.Rect)) {
		return
	}
	p.Pix[y-p.Rect.Min.Y][x-p.Rect.Min.X].Add(c)
}

// SetRGBA64 sets the pixel at (x, y) to a given color of type color.RGBA64.
func (p *LabA) SetRGBA64(x, y int, c color.RGBA64) {
	if !(image.Point{x, y}.In(p.Rect)) {
		return
	}
	if c.A == 0 {
		p.Pix[y-p.Rect.Min.Y][x-p.Rect.Min.X] = accumcolor.LabA{Tally: 1}
		return
	}
	alpha := float64(c.A)
	clr := colorful.Color{
		R: float64(c.R) / alpha,
		G: float64(c.G) / alpha,
		B: float64(c.B) / alpha,
	}
	L, a, b := clr.Lab()
	p.Pix[y-p.Rect.Min.Y][x-p.Rect.Min.X] = accumcolor.LabA{
		L:     L,
		A:     a,
		B:     b,
		Alpha: uint64(c.A >> 8),
		Tally: 1,
	}
}

// AddRGBA64 accumulates a given color of type color.RGBA64 to the pixel at (x, y).
func (p *LabA) AddRGBA64(x, y int, c color.RGBA64) {
	if !(image.Point{x, y}.In(p.Rect)) {
		return
	}
	if c.A == 0 {
		p.Pix[y-p.Rect.Min.Y][x-p.Rect.Min.X].Add(accumcolor.LabA{Tally: 1})
		return
	}
	alpha := float64(c.A)
	clr := colorful.Color{
		R: float64(c.R) / alpha,
		G: float64(c.G) / alpha,
		B: float64(c.B) / alpha,
	}
	L, a, b := clr.Lab()
	p.Pix[y-p.Rect.Min.Y][x-p.Rect.Min.X].Add(accumcolor.LabA{
		L:     L,
		A:     a,
		B:     b,
		Alpha: uint64(c.A >> 8),
		Tally: 1,
	})
}

// SubImage returns an image representing the portion of the image p visible
// through r. The returned value shares pixels with the original image.
func (p *LabA) SubImage(r image.Rectangle) image.Image {
	r = r.Intersect(p.Rect)
	// If r1 and r2 are Rectangles, r1.Intersect(r2) is not guaranteed to
	// be inside either r1 or r2.  We explicitly check for this situation.
	if r.Empty() {
		return &LabA{}
	}

	// Create new rows, each containing a subset of the original columns.
	wd, ht := r.Dx(), r.Dy()
	xOfs := r.Min.X - p.Rect.Min.X
	yOfs := r.Min.Y - p.Rect.Min.Y
	pixels := make([][]accumcolor.LabA, ht)
	for r := range pixels {
		pixels[r] = p.Pix[yOfs+r][xOfs : xOfs+wd]
	}
	return &LabA{
		Pix:  pixels,
		Rect: r,
	}
}
