package maps

import (
	"fmt"
	"path"
	"strings"

	tiled "github.com/lafriks/go-tiled"
	"github.com/oakmound/oak/v2/dlog"
	"github.com/oakmound/oak/v2/render"
	"github.com/pkg/errors"
)

// Level encapsulates the oak idea of the tiles that make up a level.
// For the sake of the example this is living for now in the example.
// This may later implement an interface which lives in oak/entites/x but for now lives here.
type Level [][][]Tile

// LoadLevelFromTMX will attempt to use the underlying library to load a TMX definition.
// It will then take the resultant level and return an oaken version :).
func LoadLevelFromTMX(fileToLoad string) (levelPointer *Level, tmxMap *tiled.Map, err error) {

	dlog.Info("attempting to load and convert from tmx: ", fileToLoad)
	// Parse a .tmx file with the go-tiled library.
	// The library has some pretty big gaps in what it handles such as compression types,
	// but seems to still be updated and is a good start for now.
	gameMap, err := tiled.LoadFromFile(fileToLoad)
	if err != nil {
		return levelPointer, gameMap, errors.Wrap(err, "error parsing map")
	}
	tmxMap = gameMap
	// for now we will only support the same subset as go-tiled does for rendering
	// see their function for a new renderer for where they perform this check.
	if !(gameMap.Orientation == "orthogonal") {
		return levelPointer, tmxMap, errors.Errorf("unsupported orientation %s, currently only support 'orthogonal'", gameMap.RenderOrder)
	}
	// see their function RenderLayer for where they perform this check.
	if !(gameMap.RenderOrder == "" || gameMap.RenderOrder == "right-down") {
		return levelPointer, tmxMap, errors.Errorf("unsupported renderOrder %s, currently only support ''/'right-down'", gameMap.RenderOrder)
	}

	// Given a mismatch of locationality its actually important for us to now
	// get the actual directory that the file we loaded is in because of the possiblity of relative links...

	tmxDir := path.Dir(strings.Replace(fileToLoad, "\\", "/", -1))
	dlog.Info("attempting to access all renderables from ", tmxDir)

	// spew.Config.MaxDepth = 1
	// spew.Dump(gameMap)

	// Now that we have the tmx map time to create our structure for its storage
	dlog.Info("TMX has loaded, starting on the conversion.")

	// Sneaky stuff here to preload assets for conversion?

	// Preload included tilesets though note we dont take into account margin in addition to spacing at this time.
	tilesetsByName := map[string]*render.Sheet{}
	// If you need that in your assets for now you can rely on the inbuilt functions of go-tiled.
	// However there maybe another function that gets added to oak for this functionality.
	for _, tSet := range gameMap.Tilesets {
		dlog.Warn("loading tmx tileset:", tSet.Name)
		tilesetsByName[tSet.Name], err = render.LoadSheet(tmxDir, tSet.Image.Source, tSet.TileWidth, tSet.TileHeight, tSet.Spacing)
		if err != nil {
			dlog.Warn(err)
		}
	}

	// Naively assume for now that the user wants all layers to be converted
	l := make([][][]Tile, len(gameMap.Layers))
	for z, tmxLayer := range gameMap.Layers {
		dlog.Info("loading tmxLayer:", tmxLayer.Name)
		// construct the arrays to put tiles into.
		l[z] = make([][]Tile, gameMap.Width)
		for x := 0; x < gameMap.Width; x++ {
			l[z][x] = make([]Tile, gameMap.Height)
			for y := 0; y < gameMap.Height; y++ {
				l[z][x][y] = &Basic{}
			}
		}
		// For now to have parity with the rendering style recommended in go-tiled lets just do this.
		// #ignore the invisible
		if !tmxLayer.Visible {
			continue
		}

		// for now since we are only supporting right down we didnt pull this into a function
		// TODO: pull this out to support different orientations and perhaps map types
		tIndex := 0
		for y := 0; y < gameMap.Height; y++ {
			for x := 0; x < gameMap.Width; x++ {
				sourceTile := tmxLayer.Tiles[tIndex]
				tIndex++
				if sourceTile.IsNil() {
					continue
				}
				tSet, ok := tilesetsByName[sourceTile.Tileset.Name]
				if !ok {
					return levelPointer, tmxMap, errors.Errorf("unloaded tileset: %s", sourceTile.Tileset.Name)
				}
				// THIS is why the load needs to be refactored
				// offsetIDs := int(sourceTile.ID - sourceTile.Tileset.FirstGID)
				offsetIDs := int(sourceTile.ID)

				subX := (offsetIDs) % sourceTile.Tileset.Columns
				subY := (offsetIDs / sourceTile.Tileset.Columns)
				if subX > len(*tSet) {
					ss := tSet.ToSprites()
					dlog.Warn(fmt.Sprintf("You are probably going to regret this.... %d %d requested from maxes of %d, %d from the tile %d",
						subX, subY, len(*tSet), len(ss[0]), sourceTile.ID))
				}
				tImg := tSet.SubSprite(subX, subY)

				// We would copy but subsprite is a bad method that dont cache
				l[z][x][y].AssignRenderable(tImg)
				l[z][x][y].ShiftPos(float64(x*gameMap.TileWidth), float64(y*gameMap.TileHeight)).SetLayer(((z + 1) * 1)) // for now its just overlay on the z layer bu't thats simplistic.

				l[z][x][y].Draw()
			}
		}

	}

	return
}
