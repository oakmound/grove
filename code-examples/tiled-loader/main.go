package main

import (
	"path/filepath"

	"github.com/implausiblyfun/oakgrove/code-examples/tiled-loader/maps"
	oak "github.com/oakmound/oak/v2"
	"github.com/oakmound/oak/v2/dlog"
	"github.com/oakmound/oak/v2/entities"
	"github.com/oakmound/oak/v2/entities/x/move"
	"github.com/oakmound/oak/v2/physics"
	"github.com/oakmound/oak/v2/render"
	"github.com/oakmound/oak/v2/render/mod"
	"github.com/oakmound/oak/v2/scene"
)

// just for fun set the name of the example for pathing.
const exampleName = "tiled-loader"

func main() {

	// Note: you can also set the SetupConfig in a single file.

	// Set the asset path to account for the special location due to the mixed licensing.
	oak.SetupConfig.Assets.AssetPath = filepath.Join("..", "..", "assets", exampleName)
	// This indicates to oak to automatically open and load image and audio
	// files local to the project before starting any scene.
	oak.SetupConfig.BatchLoad = true
	oak.SetupConfig.Debug.Level = "INFO"

	oak.Add("tiled-walk", func(string, interface{}) {

		// Get the map info.
		// Given that this is an example we will reload from file each time the scene starts.
		ourMap := filepath.Join(oak.SetupConfig.Assets.AssetPath, "maps", "map1.tmx")
		_, info, err := maps.LoadLevelFromTMX(ourMap)
		if err != nil {
			panic(err.Error())
		}
		dlog.Info("Sizing is ", info.Width*info.TileWidth, info.Height)
		oak.SetViewportBounds(0, 0, info.Width*info.TileWidth, info.Height*info.TileHeight)

		// Load a sheet to draw our character from. Art thanks to Fenrir!
		pSheet, err := render.GetSheet(filepath.Join("", "16x32", "droid_1_sheet.png"))
		dlog.ErrorCheck(err)

		// Load the animation sequences and ingore the errors cause we are lazy.
		hoverHZ, _ := render.NewSheetSequence(pSheet, 6, 0, 0, 1, 0, 2, 0, 1, 0)
		hoverUP, _ := render.NewSheetSequence(pSheet, 6, 0, 2, 1, 2, 2, 2, 1, 2)
		hoverDN, _ := render.NewSheetSequence(pSheet, 6, 0, 1, 1, 1, 2, 1, 1, 1)

		playerR := render.NewSwitch("left", map[string]render.Modifiable{
			"left":  hoverHZ.Copy().Modify(mod.FlipX),
			"right": hoverHZ.Copy(),
			"up":    hoverUP.Copy(),
			"down":  hoverDN.Copy(),
		})
		if err != nil {
			dlog.Error(err)
		}
		char := entities.NewMoving(100, 100, 32, 32,
			playerR,
			nil, 0, 0)
		char.Speed = physics.NewVector(3, 3)
		char.Bind(func(int, interface{}) int {
			// utility for moving char via wasd
			move.WASD(char)
			// the above move has now set the delta for us to key off

			toSet := playerR.Get()
			switch horiz := char.Delta.X(); {
			case horiz > 0:
				toSet = "right"
			case horiz < 0:
				toSet = "left"
			default:
				if char.Delta.Y() > 0 {
					toSet = "down"
				} else if char.Delta.Y() < 0 {
					toSet = "up"
				}

			}
			playerR.Set(toSet)
			oak.SetScreen(
				int(char.R.X())-oak.ScreenWidth/2,
				int(char.R.Y())-oak.ScreenHeight/2,
			)
			return 0
		}, "EnterFrame")
		render.Draw(char.R, 2, 2)

	}, func() bool {
		return true
	}, func() (string, *scene.Result) {
		return "tiled-walk", nil
	})

	render.SetDrawStack(
		render.NewCompositeR(),
		render.NewHeap(false),
		render.NewHeap(false),
		render.NewDrawFPS(),
		render.NewLogicFPS(),
	)

	oak.Init("tiled-walk")
}
