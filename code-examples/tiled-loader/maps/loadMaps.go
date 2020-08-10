package maps

import (
	"github.com/davecgh/go-spew/spew"
	tiled "github.com/lafriks/go-tiled"
	"github.com/lafriks/go-tiled/render"
	"github.com/pkg/errors"
	// "github.com/lafriks/go-tiled/render"
)

const mapPath = "../../assets/tiled-loader/maps/map1.tmx" // Path to your Tiled Map.

func Load() error {
	// Parse .tmx file.
	gameMap, err := tiled.LoadFromFile(mapPath)

	if err != nil {
		return errors.Wrap(err, "error parsing map")
	}

	spew.Print(gameMap)

	// You can also render the map to an in-memory image for direct
	// use with the default Renderer, or by making your own.
	renderer, _ := render.NewRenderer(gameMap)

	// Render just layer 0 to the Renderer.
	renderer.RenderLayer(0)

	// Get a reference to the Renderer's output, an image.NRGBA struct.
	//img := renderer.Result

	// Clear the render result after copying the output if separation of
	// layers is desired.
	// renderer.Clear()

	// And so on. You can also export the image to a file by using the
	// Renderer's Save functions.

	return nil
}
