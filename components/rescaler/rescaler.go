package rescaler

import (
	"github.com/oakmound/oak/v4/alg/floatgeom"
	"github.com/oakmound/oak/v4/alg/intgeom"
	"github.com/oakmound/oak/v4/entities/x/btn"
	"github.com/oakmound/oak/v4/render"
	"github.com/oakmound/oak/v4/render/mod"
	"github.com/oakmound/oak/v4/scene"
)

// A Rescaler knows the resolution a scene is designed for, and can make simple
// adjustments to scale and positioning of elements in the scene if the window
// is scaled differently from the target resolution at scene start. The rescaler
// does not handle a window's size changing mid-scene.
type Rescaler struct {
	ctx              *scene.Context
	targetResolution intgeom.Point2
	wRatio, hRatio   float64
}

// New should take in the resolution the scene's elements are all designed for.
func New(ctx *scene.Context, targetResolution intgeom.Point2) *Rescaler {
	rsc := &Rescaler{
		ctx:              ctx,
		targetResolution: targetResolution,
	}
	w, h := rsc.ctx.Window.Width(), rsc.ctx.Window.Height()
	rsc.wRatio = float64(w) / float64(rsc.targetResolution.X())
	rsc.hRatio = float64(h) / float64(rsc.targetResolution.Y())
	return rsc
}

func (rsc Rescaler) Position(target floatgeom.Point2) (float64, float64) {
	return target.X() * rsc.wRatio, target.Y() * rsc.hRatio
}

func (rsc Rescaler) SetPosition(rnd render.Renderable, target floatgeom.Point2) {
	rnd.SetPos(rsc.Position(target))
}

// Scale will scale a modifiable
func (rsc Rescaler) Scale(m render.Modifiable) {
	m.Modify(mod.Scale(rsc.wRatio, rsc.hRatio))
}

func (rsc Rescaler) Height(h float64) float64 {
	return h * rsc.hRatio
}

func (rsc Rescaler) Width(w float64) float64 {
	return w * rsc.wRatio
}

func (rsc Rescaler) Dims(dims floatgeom.Point2) floatgeom.Point2 {
	dims[0] *= rsc.hRatio
	dims[1] *= rsc.wRatio
	return dims
}

func (rs Rescaler) Draw(r render.Renderable, layers ...int) {
	if m, ok := r.(render.Modifiable); ok {
		rs.Scale(m)
	}
	rs.ctx.DrawStack.Draw(r, layers...)
}

// BtnOption can be used with an x/btn option set to apply most operations a
// rescaler provides where meaningful.
func BtnOption(rsc *Rescaler) btn.Option {
	return func(g btn.Generator) btn.Generator {
		if g.W != 0 {
			g.W = rsc.Width(g.W)
		}
		if g.H != 0 {
			g.H = rsc.Height(g.H)
		}
		if g.X != 0 || g.Y != 0 {
			g.X, g.Y = rsc.Position(floatgeom.Point2{g.X, g.Y})
		}
		if g.TxtX != 0 || g.TxtY != 0 {
			g.TxtX, g.TxtY = rsc.Position(floatgeom.Point2{g.TxtX, g.TxtY})
		}
		return g
	}
}
