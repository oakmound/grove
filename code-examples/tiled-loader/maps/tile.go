package maps

import "github.com/oakmound/oak/v2/render"

type Tile interface {
	Draw() error
	AssignRenderable(render.Renderable) Tile
	GetR() render.Renderable
	SetPos(x, y float64) Tile
	ShiftPos(x, y float64) Tile
	SetLayer(l int) Tile
	IsNil() bool
	Copy() Tile
}

// ID of the type of thing the tile is.
type ID int

const (
	Nothing ID = iota
	Floor
	Wall
	Pain
)
