package main

import (
	"fmt"
	"image/color"
	"os"
	"time"

	"github.com/oakmound/grove/components/fonthelper"
	"github.com/oakmound/grove/components/keyhint"
	"github.com/oakmound/grove/components/sound"
	"github.com/oakmound/grove/components/textfit"
	"github.com/oakmound/grove/components/textinput"
	"github.com/oakmound/grove/components/textqueue"
	"golang.org/x/image/colornames"

	"github.com/oakmound/oak/v4"
	"github.com/oakmound/oak/v4/alg/floatgeom"
	"github.com/oakmound/oak/v4/audio"
	"github.com/oakmound/oak/v4/collision"
	"github.com/oakmound/oak/v4/debugtools"
	"github.com/oakmound/oak/v4/entities"
	"github.com/oakmound/oak/v4/event"
	"github.com/oakmound/oak/v4/key"
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
		render.Draw(debugtools.NewThickRTree(ctx, ctx.CollisionTree, 5), 2, 3)
		render.Draw(debugtools.NewThickColoredRTree(ctx, ctx.MouseTree, 5,
			map[collision.Label]color.RGBA{collision.Label(0): {200, 100, 100, 255}}), 2, 3)
		basePosI := ctx.Window.Bounds().DivConst(2)
		basePos := floatgeom.Point2{float64(basePosI.X()), float64(basePosI.Y())}
		err := audio.InitDefault()
		if err != nil {
			fmt.Println("init failed:", err)
			os.Exit(1)
		}
		sound.Init(1, 1, 1)
		mainBar := sound.NewBar(ctx, sound.KindMaster, basePos, 100, 40)
		ctx.Draw(mainBar.Renderable)

		event.Bind(ctx, mouse.PressOn, mainBar, func(bar *entities.Entity, me *mouse.Event) event.Response {
			sound.SetMasterVolume(ctx, solidProgress(bar, *me))
			return 0
		})

		event.Bind(ctx, key.Down(key.R), mainBar, func(bar *entities.Entity, kv key.Event) event.Response {
			sound.SetMasterVolume(ctx, .5)
			return 0
		})

		fnt, _ := render.DefFontGenerator.RegenerateWith(
			fonthelper.WithSizeAndColor(10, colornames.Purple))
		hint := keyhint.Options{Key: "R", Height: 20, Font: fnt, Color: colornames.Antiquewhite}.Generate()
		keyhint.AlignHintToButton(hint, mainBar)
		ctx.Draw(hint)

		kFnt, _ := render.DefFontGenerator.RegenerateWith(
			fonthelper.WithSizeAndColor(15, colornames.Orange))
		keyQueue := textqueue.New(
			ctx, []event.UnsafeEventID{key.AnyDown.UnsafeEventID, mouse.Press.UnsafeEventID},
			floatgeom.Point2{50, 50}, 0,
			kFnt, 4*time.Second,
		)

		ctx.DrawStack.Draw(keyQueue, 0)

		wrapFnt, err := render.DefFontGenerator.RegenerateWith(
			fonthelper.WithSizeAndColor(15, colornames.Palevioletred))
		wrappedText, err := textfit.New(
			textfit.String(
				"This scene has a collection of potentially useful things for management of options "+
					"and things like debug scenes or other ways to display event info "+
					"This is a smattering of random grove components and is prone to change in the future. "+
					"For example this is a instance of textfit and wrapping of text. ",
			),
			textfit.Dimensions(floatgeom.Point2{200, 200}),
			textfit.Font(wrapFnt),
			textfit.Inset(15, 15),
			textfit.WithBreakStyle(textfit.BreakStyleWord),
			textfit.MinSize(5),
			textfit.MaxSize(22),
		)

		entities.New(ctx,
			entities.WithPosition(floatgeom.Point2{basePos.X(), 0}),
			entities.WithRenderable(wrappedText),
			entities.WithDrawLayers([]int{1, 2}),
		)

		textinput.New(ctx,
			textinput.WithPos(basePos.X(), basePos.Y()+100),
			textinput.WithDims(200, 22),
			textinput.WithStr("Input example"),
			//textinput.WithEntityOptions(entities.WithDrawLayers([]int{1, 2}))
		)
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
