// This file presents a few examples of accumimage.

package accumimage_test

import (
	"image"
	"image/png"
	"log"
	"os"

	"github.com/spakin/accumimage"
)

// readImage reads an image from a named file.
func readImage(fn string) (image.Image, error) {
	r, err := os.Open(fn)
	if err != nil {
		return nil, err
	}
	img, _, err := image.Decode(r)
	if err != nil {
		return nil, err
	}
	return img, nil
}

// writeImage writes an image to a named file in PNG format.
func writeImage(fn string, img image.Image) error {
	w, err := os.Create(fn)
	if err != nil {
		return err
	}
	err = png.Encode(w, img)
	return err
}

// Scale down an arbitrary image (img) to given dimensions (newBnds), averaging
// colors that map to the same target pixel.  The result is smoother than if
// newImg.Set were used instead of newImg.Add.
func Example() {
	// Read an image from a file.
	if len(os.Args) != 2 {
		return
	}
	img, err := readImage(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	// Choose (arbitrarily) to scale to 1/3 of the size in each of x and y.
	bnds := img.Bounds()
	wd, ht := bnds.Dx(), bnds.Dy()
	newBnds := image.Rect(0, 0, wd/3, ht/3)

	// Create an AccumNRGBA image.
	newImg := accumimage.NewAccumNRGBA(newBnds)
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

	// Write the downscaled image to a file.
	err = writeImage("image.png", newImg)
	if err != nil {
		log.Fatal(err)
	}
}

// Blend two images to produce a third image.
func ExampleAccumNRGBA_Add() {
	// Read two images from files.
	if len(os.Args) != 3 {
		return
	}
	imgA, err := readImage(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	imgB, err := readImage(os.Args[2])
	if err != nil {
		log.Fatal(err)
	}

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

	// Write the blended image to a file.
	err = writeImage("image.png", imgC)
	if err != nil {
		log.Fatal(err)
	}
}
