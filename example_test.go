// This file presents a few examples of accumimage.

package accumimage_test

import (
	"image"

	"github.com/spakin/accumimage"
)

var img, imgA, imgB image.Image // Images used by examples
var newBnds image.Rectangle     // Rectangle used by examples

// Scale down an arbitrary image (img) to given dimensions (newBnds), averaging
// colors that map to the same target pixel.  The result is smoother than if
// newImg.Set were used instead of newImg.Add.
func Example() {
	// Create an AccumNRGBA image.
	newImg := accumimage.NewAccumNRGBA(newBnds)

	// Acquire the dimensions of both the old and new images.
	bnds := img.Bounds()
	wd, ht := bnds.Dx(), bnds.Dy()
	nwd, nht := newBnds.Dx(), newBnds.Dy()

	// Copy the image pixel-by-pixel, from (x, y) in the input image to
	// (nx, ny) in the output image.
	for y := bnds.Min.Y; y < bnds.Max.Y; y++ {
		ny := (nht*(y-bnds.Min.Y))/ht + newBnds.Min.Y
		for x := bnds.Min.X; x < bnds.Max.X; x++ {
			nx := (nwd*(x-bnds.Min.X))/wd + newBnds.Min.X
			c := img.At(x, y)
			newImg.Add(nx, ny, c) // Accumulate multiple colors into a single pixel.
		}
	}
}

// Blend two images to produce a third image.
func ExampleAccumNRGBA_Add() {
	// Create an AccumNRGBA image imgC that's large enough to hold the
	// overlap of input images imgA and imgB.
	bndsA := imgA.Bounds()
	bndsB := imgB.Bounds()
	bndsC := bndsA.Union(bndsB)
	imgC := accumimage.NewAccumNRGBA(bndsC)

	// Blend image A with image C.
	for y := bndsA.Min.Y; y < bndsA.Max.Y; y++ {
		for x := bndsA.Min.X; x < bndsA.Max.X; x++ {
			imgC.Add(x, y, imgA.At(x, y))
		}
	}

	// Blend image B with image C.
	for y := bndsB.Min.Y; y < bndsB.Max.Y; y++ {
		for x := bndsB.Min.X; x < bndsB.Max.X; x++ {
			imgC.Add(x, y, imgB.At(x, y))
		}
	}
}
