package accumimage_test

import "github.com/spakin/accumimage"

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
	return newImg
}
