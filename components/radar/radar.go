package radar

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"sync"

	"github.com/oakmound/oak/v4/render"
)

// Point is a utility function for location
type Point struct {
	X, Y *float64
}

// Radar displays points of interest on a radar map
type Radar struct {
	render.LayeredPoint
	points        map[Point]color.Color
	pointLookup   map[int]*Point
	center        Point
	width, height int
	r             *image.RGBA
	outline       *render.Sprite
	ratio         float64
	sync.Mutex
}

var (
	centerColor = color.RGBA{255, 255, 0, 255}
)

// NewRadar creates a radar that will display at 0,0 with the given dimensions.
// The points given will be displayed on the radar relative to the center point,
// With the absolute distance reduced by the given ratio
func NewRadar(w, h int, points map[Point]color.Color, center Point, ratio float64) *Radar {
	r := new(Radar)
	r.LayeredPoint = render.NewLayeredPoint(0, 0, 0)
	r.points = points
	r.pointLookup = map[int]*Point{}
	r.width = w
	r.height = h
	r.center = center
	r.r = image.NewRGBA(image.Rect(0, 0, w, h))
	r.outline = render.NewColorBox(w, h, color.RGBA{0, 0, 125, 125})
	r.ratio = ratio
	return r
}

// SetPos sets the position of the radar on the screen
func (r *Radar) SetPos(x, y float64) {
	r.LayeredPoint.SetPos(x, y)
	r.outline.SetPos(x, y)
}

// SetOutline of the radar to be the provided outline
func (r *Radar) SetOutline(outline *render.Sprite) {
	r.outline = outline
}

// GetRGBA returns this radar's image
func (r *Radar) GetRGBA() *image.RGBA {
	return r.r
}

// Draw draws the radar at a given offset
func (r *Radar) Draw(buff draw.Image, xOff, yOff float64) {
	// Draw each point p in r.points
	// at r.X() + center.X() - p.X(), r.Y() + center.Y() - p.Y()
	// IF that value is < r.width/2, > -r.width/2, < r.height/2, > -r.height/2
	r.Lock()
	for p, c := range r.points {
		x := int((*p.X-*r.center.X)/r.ratio) + r.width/2
		y := int((*p.Y-*r.center.Y)/r.ratio) + r.height/2
		for x2 := x - 1; x2 < x+1; x2++ {
			for y2 := y - 1; y2 < y+1; y2++ {
				r.r.Set(x2, y2, c)
			}
		}
	}
	r.Unlock()
	r.r.Set(r.width/2, r.height/2, centerColor)
	r.outline.Draw(buff, xOff, yOff)
	render.DrawImage(buff, r.r, int(xOff+r.X()), int(yOff+r.Y()))

	r.r = image.NewRGBA(image.Rect(0, 0, r.width, r.height))
}

// AddPoint adds an additional point to the radar to be tracked
func (r *Radar) AddPoint(loc Point, c color.Color) {
	r.Lock()
	r.points[loc] = c
	r.Unlock()
}

// AddTrackedPoint to the radar. Enables display and later lookup by id (usually the CID of the caller).
func (r *Radar) AddTrackedPoint(loc Point, id int, c color.Color) {
	r.Lock()
	r.points[loc] = c
	r.pointLookup[id] = &loc
	r.Unlock()
}

// LookupPoint by the provided id. This only works if the point was tracked on creation.
func (r *Radar) LookupPoint(id int) (*Point, bool) {
	p, ok := r.pointLookup[id]
	return p, ok
}

// RemovePointByLookup removes a point if it is present. If it is not, it does nothing.
func (r *Radar) RemovePointByLookup(id int) error {
	p, ok := r.LookupPoint(id)
	if !ok {
		return nil
	}
	loc := *p
	r.Lock()
	defer r.Unlock()
	if _, ok := r.points[loc]; ok {
		delete(r.points, loc)
	} else {
		return fmt.Errorf("attempted to remove a radar point that did not exist")
	}
	return nil
}
