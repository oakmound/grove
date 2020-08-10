package maps

import (
	"github.com/oakmound/oak/v2/render"
)

type Basic struct {
	r  render.Renderable
	id ID
}

func (bt *Basic) Draw() error {
	if !bt.IsNil() {
		render.Draw(bt.r, 2, bt.r.GetLayer())
	}
	return nil
}

func (bt *Basic) AssignRenderable(r render.Renderable) Tile {
	bt.r = r
	return bt
}

func (bt *Basic) Id() ID {
	return bt.id
}
func (bt *Basic) GetR() render.Renderable {
	return bt.r
}

func (bt *Basic) ShiftPos(x, y float64) Tile {
	if !bt.IsNil() {
		bt.r.ShiftX(x)
		bt.r.ShiftY(y)
	}
	return bt
}

func (bt *Basic) SetPos(x, y float64) Tile {
	if !bt.IsNil() {
		bt.r.SetPos(x, y)
	}
	return bt
}

func (bt *Basic) SetLayer(l int) Tile {
	if !bt.IsNil() {
		bt.r.SetLayer(l)
	}
	return bt
}
func (bt *Basic) IsNil() bool {
	return bt == nil || bt.r == nil
}

func (bt *Basic) Copy() Tile {
	bt2 := *bt
	return &bt2
}
