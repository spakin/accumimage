accumimage and accumcolor
=========================
[![Go project version](https://badge.fury.io/go/github.com%2Fspakin%2Faccumimage%2Fv2.svg)](https://badge.fury.io/go/github.com%2Fspakin%2Faccumimage%2Fv2)
[![Go](https://github.com/spakin/accumimage/actions/workflows/go.yml/badge.svg)](https://github.com/spakin/accumimage/actions/workflows/go.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/spakin/accumimage)](https://goreportcard.com/report/github.com/spakin/accumimage/v2)
[![Go Reference](https://pkg.go.dev/badge/github.com/spakin/accumimage.svg)](https://pkg.go.dev/github.com/spakin/accumimage/v2)

Description
-----------

`accumimage` and `accumcolor` are packages for the [Go programming language](https://golang.org/) that provide the ability to accumulate and average colors.  `accumcolor.NRGBA` is the accumulating analogue of the Go standard library's [`color.NRGBA`](https://pkg.go.dev/image/color#NRGBA) (a single pixel of non-alpha-premultiplied red, green, blue, and alpha channels), and `accumimage.NRGBA` is the accumulating analogue of the Go standard library's [`image.NRGBA`](https://pkg.go.dev/image#NRGBA) (a 2-D array of NRGBA pixels).  Another type, `accumcolor.LabA`, accumulates color in the [CIE L\*a\*b\*](https://en.wikipedia.org/wiki/CIELAB_color_space) + alpha color space, in which colors blend in a more perceptually uniform manner.  `accumimage.LabA` provides a 2-D array of `accumcolor.LabA` pixels.

In addition to three color channels and an alpha channel, `accumcolor.NRGBA` and `accumcolor.LabA` contain a tally of the number of colors that have been added together.  The `Add` method adds another color to the existing one, each weighted by their tally.  The `NRGBA` method divides each channel in an `accumcolor.NRGBA` by the tally to produce an ordinary `color.NRGBA`; the `Colorful` method divides each channel in an `accumcolor.LabA` by the tally to produce a [`colorful.Color`](https://pkg.go.dev/github.com/lucasb-eyer/go-colorful#Color) from the [`go-colorful`](https://pkg.go.dev/github.com/lucasb-eyer/go-colorful) package.

`accumimage` and `accumcolor` can be useful for blending multiple overlapping images, for mapping a large number of pixels to a smaller number (e.g., when scaling an image), and for distinguishing pixels in an image that are truly empty (e.g., unvisited in an algorithm that visits all pixels) from pixels that have a color, even one that is fully transparent.


Usage
-----

Install `accumimage` and `accumcolor` with
```bash
go get github.com/spakin/accumimage/v2
```

then simply import `accumimage` and/or `accumimage/accumcolor` into your Go program:
```Go
import (
        "github.com/spakin/accumimage/v2"
        "github.com/spakin/accumimage/v2/accumcolor"
)
```

Documentation
-------------

See the [pkg.go.dev documentation for `accumimage`](https://pkg.go.dev/github.com/spakin/accumimage/v2) and [`accumcolor`](https://pkg.go.dev/github.com/spakin/accumimage/v2/accumcolor).


Author
------

[Scott Pakin](http://www.pakin.org/~scott/), *scott-accum@pakin.org*
