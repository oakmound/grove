package keyhint

import (
	"fmt"
	"image"
	"image/color"
	"math"

	"github.com/oakmound/oak/v3/alg/intgeom"
	"github.com/oakmound/oak/v3/entities/x/btn"
	"github.com/oakmound/oak/v3/entities/x/mods"
	"github.com/oakmound/oak/v3/render"
)

// A KeyHint is a small, colored renderable displaying some text. The text 
// is designed to be a key 'A,B,C,D' or controller button 'A,B,X,Y', and intended
// for use as a hint to imply that pressing that key or button will trigger something. 
type KeyHint struct {
	Options
	render.Renderable
}

type Options struct {
	Key         string
	Color       color.RGBA
	BorderColor color.RGBA
	Font        *render.Font
	Height      int
	// Todo: Fit text
	Rounded bool
}

func (o Options) setDefaults() Options {
	if o.Key == "" {
		o.Key = "A"
	}
	if o.Color == (color.RGBA{}) {
		o.Color = color.RGBA{200, 0, 0, 255}
	}
	if o.BorderColor == (color.RGBA{}) {
		o.BorderColor = color.RGBA{0, 0, 0, 255}
	}
	if o.Height == 0 {
		o.Height = 35
	}
	if o.Font == nil {
		o.Font = render.DefaultFont()
	}
	return o
}

func (o Options) Generate() *KeyHint {
	o = o.setDefaults()
	rgba := image.NewRGBA(image.Rect(0, 0, o.Height, o.Height))
	// color based on delta from center
	center := intgeom.Point2{o.Height / 2, o.Height / 2}
	for x := 0; x < o.Height; x++ {
		for y := 0; y < o.Height; y++ {
			pt := intgeom.Point2{x, y}
			distance := 0.0
			if o.Rounded {
				distance = center.Distance(pt)
			} else {
				// distance is greater of x or y delta
				xDelta := math.Abs(float64(center.X() - pt.X()))
				yDelta := math.Abs(float64(center.Y() - pt.Y()))
				distance = xDelta
				if yDelta > xDelta {
					distance = yDelta
				}
			}
			// color by distance:
			distancePercent := distance / float64(o.Height/2)
			var c color.Color = o.Color
			if distancePercent > 1 {
				continue
			} else if distancePercent > (float64(o.Height)-0.99)/float64(o.Height) {
				rC := o.BorderColor
				rC.A /= 4
				rC.R /= 4
				rC.B /= 4
				rC.G /= 4
				c = rC
			} else if distancePercent > (float64(o.Height)-1.99)/float64(o.Height) {
				rC := o.BorderColor
				rC.A /= 2
				rC.R /= 2
				rC.B /= 2
				rC.G /= 2
				c = rC
			} else if distancePercent > (float64(o.Height)-2.99)/float64(o.Height) {
				c = o.BorderColor
			} else if distancePercent > (float64(o.Height)-3.99)/float64(o.Height) {
				c = mix(o.BorderColor, mods.Darker(c, .2), .80)
			} else if distancePercent > (float64(o.Height)-4.99)/float64(o.Height) {
				c = mix(o.BorderColor, mods.Darker(c, .2), .50)
			} else if distancePercent > .85 {
				c = mods.Darker(c, .2)
			} else if distancePercent > .65 {
				c = mods.Lighter(c, distancePercent-.65)
			}
			rgba.Set(x, y, c)
		}
	}
	backing := render.NewSprite(0, 0, rgba)
	text := o.Font.NewText(o.Key, 0, 0)
	tw, th := text.GetDims()
	text.SetPos(float64(center.X()-tw/2), float64(center.Y()-th/2)-3)
	comp := render.NewCompositeR(backing, text.ToSprite())
	return &KeyHint{
		Renderable: comp,
		Options:    o,
	}
}

func mix(c1, c2 color.Color, percent float64) color.Color {
	r1, g1, b1, a1 := c1.RGBA()
	r2, g2, b2, a2 := c2.RGBA()
	return color.RGBA64{
		uint16(float64(r1)*(percent) + float64(r2)*(1-percent)),
		uint16(float64(g1)*(percent) + float64(g2)*(1-percent)),
		uint16(float64(b1)*(percent) + float64(b2)*(1-percent)),
		uint16(float64(a1)*(percent) + float64(a2)*(1-percent)),
	}
}

func AlignHintToButton(kh *KeyHint, b btn.Btn) {
	r := b.GetRenderable()
	w, _ := r.GetDims()
	kh.SetPos((r.X()+float64(w))-float64(kh.Height)*(3.0/4.0), r.Y()-(float64(kh.Height)*(1.0/4.0)))
}
