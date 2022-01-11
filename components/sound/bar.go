package sound

import (
	"fmt"
	"image/color"
	"image/draw"

	"github.com/oakmound/oak/v3/alg/floatgeom"
	"github.com/oakmound/oak/v3/alg/intgeom"
	"github.com/oakmound/oak/v3/entities"
	"github.com/oakmound/oak/v3/event"
	"github.com/oakmound/oak/v3/mouse"
	"github.com/oakmound/oak/v3/render"
	"golang.org/x/image/colornames"
)

// BarKind details which type of volume should be manipulated by the bar
type BarKind int

// Types of bars.
const (
	BarKindMaster BarKind = iota
	BarKindMusic  BarKind = iota
	BarKindSFX    BarKind = iota
)

// NewBar for setting of the sound graphically.
func NewBar(kind BarKind, pos floatgeom.Point2, w, h int) *entities.Solid {

	dl := &DashedLine{
		LayeredPoint: render.NewLayeredPoint(pos.X(), pos.Y(), 0),
		Dims:         intgeom.Point2{w, h},
		OnColor:      colornames.White,
		OffColor:     colornames.Gray,
		DashMod:      4,
		Progress:     0,
	}
	solid := entities.NewSolid(pos.X(), pos.Y(), float64(w), float64(h), dl, mouse.DefaultTree, 0)
	var eventName string
	switch kind {
	case BarKindMaster:
		eventName = EventMasterVolumeChanged
		dl.Progress = volume
	case BarKindMusic:
		eventName = EventMusicVolumeChanged
		dl.Progress = musicVolume
	case BarKindSFX:
		eventName = EventSFXVolumeChanged
		dl.Progress = sfxVolume
	}
	solid.CID.Bind(eventName, func(id event.CID, payload interface{}) int {
		newVal, ok := payload.(float64)
		if !ok {
			fmt.Println("expected float progress arg to bar change binding")
			return 0
		}
		solid, ok := event.GetEntity(id).(*entities.Solid)
		if !ok {
			fmt.Println("expected doodad entity in bar change binding")
			return 0
		}
		dl, ok := solid.R.(*DashedLine)
		if !ok {
			fmt.Println("expected renderable as DashedLine in bar change binding")
			return 0
		}
		dl.Progress = newVal
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
