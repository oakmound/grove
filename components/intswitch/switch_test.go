package intswitch_test

import (
	"github.com/oakmound/grove/components/intswitch"
	"github.com/oakmound/oak/v4/render"

	"image/color"
	"testing"
)

func TestNew(t *testing.T) {
	sw := intswitch.New(0, map[int]render.Modifiable{
		0: render.NewColorBox(1, 1, color.RGBA{255, 0, 0, 255}),
	})
	if sw == nil {
		t.Fatal("expected non-nil switch")
	}
}
