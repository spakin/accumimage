/*
Package accumimage provides support for images whose colors can be
accumulated and averaged.  It is based on the color types defined in
accumimage/accumcolor.  The core data types that accumimage defines
are AccumNRGBA and AccumLabA, which implement the image.Image
interface as well as most of the standard set of methods provided by
the image package's image types.  (AccumLabA lacks PixOffset.)  In
addition, each Accum____.Set* method has a corresponding
Accum____.Add* method, which adds color to a pixel rather than
replacing the pixel's color with a given color.
*/
package accumimage
