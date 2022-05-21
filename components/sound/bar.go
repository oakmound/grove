package sound

import (
	"image/color"
	"image/draw"

	"github.com/oakmound/oak/v4/alg/floatgeom"
	"github.com/oakmound/oak/v4/alg/intgeom"
	"github.com/oakmound/oak/v4/entities"
	"github.com/oakmound/oak/v4/event"
	"github.com/oakmound/oak/v4/render"
	"github.com/oakmound/oak/v4/scene"
	"golang.org/x/image/colornames"
)

// VolumeKind details which type of volume should be manipulated by the bar
type VolumeKind int

// Types of bars.
const (
	KindMaster VolumeKind = iota
	KindMusic  VolumeKind = iota
	KindSFX    VolumeKind = iota
)

// NewBar for setting of the sound graphically.
func NewBar(ctx *scene.Context, kind VolumeKind, pos floatgeom.Point2, w, h int) *entities.Entity {

	dl := &DashedLine{
		LayeredPoint: render.NewLayeredPoint(pos.X(), pos.Y(), 0),
		Dims:         intgeom.Point2{w, h},
		OnColor:      colornames.White,
		OffColor:     colornames.Gray,
		DashMod:      4,
		Progress:     0,
	}

	solid := entities.New(ctx,
		entities.WithRenderable(dl),
		entities.WithRect(floatgeom.NewRect2WH(pos.X(), pos.Y(), 32, 32)),
	)

	// solid := entities.NewSolid(pos.X(), pos.Y(), float64(w), float64(h), dl, mouse.DefaultTree, 0)
	// var eventName string
	switch kind {
	case KindMaster:
		// eventName = EventMasterVolumeChanged
		dl.Progress = volume
	case KindMusic:
		// eventName = EventMusicVolumeChanged
		dl.Progress = musicVolume
	case KindSFX:
		// eventName = EventSFXVolumeChanged
		dl.Progress = sfxVolume
	}

	event.Bind(ctx, EventVolumeChange, solid, func(bar *entities.Entity, change VolumeChangePayload) event.Response {
		if change.Kind != kind {
			return 0
		}
		dl.Progress = change.NewVolume
		return 0
	})

	return solid
}

// DashedLine to display the current value on the bar.
type DashedLine struct {
	render.LayeredPoint
	Dims     intgeom.Point2
	OnColor  color.RGBA
	OffColor color.RGBA
	DashMod  int
	Progress float64
}

// GetDims for the line.
func (dl *DashedLine) GetDims() (int, int) {
	return dl.Dims.X(), dl.Dims.Y()
}

// Draw the dashed line.
func (dl *DashedLine) Draw(buff draw.Image, xOff, yOff float64) {
	shouldDash := false
	wf := float64(dl.Dims.X())
	hf := float64(dl.Dims.Y())
	clr := dl.OnColor
	for i := 0.0; i < wf; i++ {
		if int(i)%dl.DashMod == 0 {
			shouldDash = !shouldDash
		}
		if !shouldDash {
			if i/wf > dl.Progress {
				clr = dl.OffColor
			}
			x := i + dl.X()
			for y := dl.Y(); y < dl.Y()+hf; y++ {
				buff.Set(int(x+xOff), int(y+yOff), clr)
			}
		}
	}
}
