accumimage and accumcolor
=========================

Description
-----------

`accumimage` and `accumcolor` are packages for the [Go programming language](https://golang.org/) that provide the ability to accumulate and average colors.  `accumcolor.AccumNRGBA` is the accumulating analogue of the Go standard library's [`color.NRGBA`](https://pkg.go.dev/image/color#NRGBA) (a single pixel of non-alpha-premultiplied red, green, blue, and alpha channels), and `accumimage.AccumNRGBA` is the accumulating analogue of the Go standard library's [`image.NRGBA`](https://pkg.go.dev/image#NRGBA) (a 2-D array of NRGBA pixels).

In addition to three color channels and an alpha channel, an `accumcolor.AccumNRGBA` contains a tally of the number of colors that have been added together.  The `Add` method adds another color to the existing one, both weighted by their tally.  The `NRGBA` method divides each channel by the tally to produce an ordinary `color.NRGBA`.

`accumimage` and `accumcolor` can be useful for blending multiple overlapping images, for mapping a large number of pixels to a smaller number (e.g., when scaling an image), and for distinguishing pixels in an image that are truly empty (e.g., unvisited in an algorithm that visits all pixels) from pixels that have a color, even one that is fully transparent.


Usage
-----

Assuming you are using the [Go module system](https://go.dev/blog/using-go-modules), simply import `accumimage` and `accumimage/accumcolor`:

```Go
import (
        "github.com/spakin/accumimage"
        "github.com/spakin/accumimage/accumcolor"
)
```
Then run `go mod tidy` to download and install the package files.

Documentation
-------------

See the [pkg.go.dev documentation for `accumimage`](https://pkg.go.dev/github.com/spakin/accumimage) and [`accumcolor`](https://pkg.go.dev/github.com/spakin/accumimage/accumcolor).


Author
------

[Scott Pakin](http://www.pakin.org/~scott/), *scott-accum@pakin.org*
