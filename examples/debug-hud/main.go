package main

import (
	"fmt"
	"os"

	s "github.com/oakmound/grove/components/sound"
	"github.com/oakmound/oak/v4"
	"github.com/oakmound/oak/v4/alg/floatgeom"
	"github.com/oakmound/oak/v4/audio"
	"github.com/oakmound/oak/v4/entities"
	"github.com/oakmound/oak/v4/event"
	"github.com/oakmound/oak/v4/mouse"
	"github.com/oakmound/oak/v4/render"
	"github.com/oakmound/oak/v4/scene"
)

func main() {
	render.SetDrawStack(
		render.NewCompositeR(),
		render.NewDynamicHeap(),
		render.NewStaticHeap(),
	)
	oak.AddScene("demo", scene.Scene{Start: func(ctx *scene.Context) {

		render.Draw(render.NewDrawFPS(0.03, nil, 10, 10))

		basePosI := ctx.Window.Bounds().DivConst(2)
		basePos := floatgeom.Point2{float64(basePosI.X()), float64(basePosI.Y())}
		err := audio.InitDefault()
		if err != nil {
			fmt.Println("init failed:", err)
			os.Exit(1)
		}
		s.Init(1, 1, 1)
		mainBar := s.NewBar(ctx, s.KindMaster, basePos, 100, 20)
		ctx.Draw(mainBar.Renderable)

		event.Bind(ctx, mouse.Click, mainBar, func(bar *entities.Entity, me *mouse.Event) event.Response {
			s.SetMasterVolume(ctx, solidProgress(bar, *me))

			return 0
		})

	}})

	oak.Init("demo")
}

func solidProgress(solid *entities.Entity, mev mouse.Event) float64 {
	const border = 10

	w, _ := solid.Renderable.GetDims()
	progress := (mev.X() - (solid.X() + border)) / (float64(w) - border*2)
	if progress < 0 {
		return 0
	}
	if progress > 1 {
		return 1
	}
	return progress
}
