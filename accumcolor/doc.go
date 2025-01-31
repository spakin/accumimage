/*
Package accumcolor provides support for colors that can be accumulated
and averaged.

An AccumNRGBA is analogous to a color.NRGBA but additionally supports adding
colors together.  It maintains a tally of the total number of colors that have
been accumulated.  The Add method accumulates more color onto an existing
AccumNRGBA.  The NRGBA method returns the average color of the entire
accumulation as a color.NRGBA.

An AccumLabA provides similar functionality to AccumNRGBA but stores colors
in CIE L*a*b* + alpha channels.  AccumLabA thereby supports a more
perceptually uniform color space and more natural results when averaging
colors.
*/
package accumcolor
