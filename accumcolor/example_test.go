// This file presents examples of accumcolor usage.

package accumcolor_test

import (
	"fmt"
	"image/color"

	"github.com/spakin/accumimage/accumcolor"
)

// Provide a variety of examples of valid and invalid NRGBA values.
func ExampleNRGBA_Valid() {
	var c accumcolor.NRGBA
	fmt.Printf("%v --> %v\n", c, c.Valid())
	c = accumcolor.NRGBA{R: 0, G: 0, B: 0, A: 0, Tally: 1}
	fmt.Printf("%v --> %v\n", c, c.Valid())
	c = accumcolor.NRGBA{R: 255, G: 255, B: 255, A: 255, Tally: 1}
	fmt.Printf("%v --> %v\n", c, c.Valid())
	c = accumcolor.NRGBA{R: 255, G: 255, B: 255, A: 255, Tally: 0}
	fmt.Printf("%v --> %v\n", c, c.Valid())
	c = accumcolor.NRGBA{R: 2550, G: 1280, B: 640, A: 2550, Tally: 1}
	fmt.Printf("%v --> %v\n", c, c.Valid())
	c = accumcolor.NRGBA{R: 2550, G: 1280, B: 640, A: 2550, Tally: 10}
	fmt.Printf("%v --> %v\n", c, c.Valid())
	// Output:
	// {0 0 0 0 0} --> true
	// {0 0 0 0 1} --> true
	// {255 255 255 255 1} --> true
	// {255 255 255 255 0} --> false
	// {2550 1280 640 2550 1} --> false
	// {2550 1280 640 2550 10} --> true
}

// Show how to average multiple NRGBA colors to produce a new NRGBA color.
func ExampleNRGBA_NRGBA() {
	c1 := color.NRGBA{R: 150, G: 100, B: 40, A: 255}
	c2 := color.NRGBA{R: 50, G: 40, B: 80, A: 255}
	var c accumcolor.NRGBA
	c.Add(c1)
	c.Add(c2)
	fmt.Printf("The average of %v and %v is %v.\n", c1, c2, c.NRGBA())
	c3 := color.NRGBA{R: 75, G: 32, B: 220, A: 255}
	c.Add(c3)
	fmt.Printf("The average of %v and %v and %v is %v.\n", c1, c2, c3, c.NRGBA())
	// Output:
	// The average of {150 100 40 255} and {50 40 80 255} is {100 70 60 255}.
	// The average of {150 100 40 255} and {50 40 80 255} and {75 32 220 255} is {91 57 113 255}.
}
