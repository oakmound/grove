package intswitch

import (
	"image"
	"image/draw"
	"sync"
	"strconv"

	"github.com/oakmound/oak/v3/event"
	"github.com/oakmound/oak/v3/oakerr"
	"github.com/oakmound/oak/v3/physics"
	"github.com/oakmound/oak/v3/render"
	"github.com/oakmound/oak/v3/render/mod"
)

var _ render.Modifiable = &Switch{}

// The Switch type is a duplicate of oak's string switch type with int-keys
// instead of string-keys. Consider using it with bitflags!
type Switch struct {
	render.LayeredPoint
	subRenderables map[int]render.Modifiable
	curRenderable  int
	lock           sync.RWMutex
}

// New creates a new Switch from a map of values to modifiables
func New(start int, m map[int]render.Modifiable) *Switch {
	return &Switch{
		LayeredPoint:   render.NewLayeredPoint(0, 0, 0),
		subRenderables: m,
		curRenderable:  start,
		lock:           sync.RWMutex{},
	}
}

// Add makes a new entry in the Switch's map. If the key already
// existed, it will be overwritten and an error will be returned.
func (c *Switch) Add(k int, v render.Modifiable) (err error) {
	if _, ok := c.subRenderables[k]; ok {
		err = oakerr.ExistingElement{
			InputName:   "k",
			InputType:   "int",
			Overwritten: true,
		}
	}
	c.lock.Lock()
	c.subRenderables[k] = v
	c.lock.Unlock()
	return err
}

// Set sets the current renderable to the one specified
func (c *Switch) Set(k int) error {
	c.lock.RLock()
	if _, ok := c.subRenderables[k]; !ok {
		return oakerr.NotFound{InputName: "k:" + strconv.Itoa(k)}
	}
	c.lock.RUnlock()
	c.curRenderable = k
	return nil
}

// GetSub returns a keyed Modifiable from this Switch's map
func (c *Switch) GetSub(s int) render.Modifiable {
	c.lock.RLock()
	m := c.subRenderables[s]
	c.lock.RUnlock()
	return m
}

// Get returns the Switch's current key
func (c *Switch) Get() int {
	return c.curRenderable
}

// SetOffsets sets the logical offset for the specified key
func (c *Switch) SetOffsets(k int, offsets physics.Vector) {
	c.lock.RLock()
	if r, ok := c.subRenderables[k]; ok {
		r.SetPos(offsets.X(), offsets.Y())
	}
	c.lock.RUnlock()
}

// Copy creates a copy of the Switch
func (c *Switch) Copy() render.Modifiable {
	newC := new(Switch)
	newC.LayeredPoint = c.LayeredPoint.Copy()
	newSubRenderables := make(map[int]render.Modifiable)
	c.lock.RLock()
	for k, v := range c.subRenderables {
		newSubRenderables[k] = v.Copy()
	}
	c.lock.RUnlock()
	newC.subRenderables = newSubRenderables
	newC.curRenderable = c.curRenderable
	newC.lock = sync.RWMutex{}
	return newC
}

//GetRGBA returns the current renderables rgba
func (c *Switch) GetRGBA() *image.RGBA {
	c.lock.RLock()
	rgba := c.subRenderables[c.curRenderable].GetRGBA()
	c.lock.RUnlock()
	return rgba
}

// Modify performs the input modifications on all elements of the Switch
func (c *Switch) Modify(ms ...mod.Mod) render.Modifiable {
	c.lock.RLock()
	for _, r := range c.subRenderables {
		r.Modify(ms...)
	}
	c.lock.RUnlock()
	return c
}

// Filter filters all elements of the Switch with fs
func (c *Switch) Filter(fs ...mod.Filter) {
	c.lock.RLock()
	for _, r := range c.subRenderables {
		r.Filter(fs...)
	}
	c.lock.RUnlock()
}

//Draw draws the Switch at an offset from its logical location
func (c *Switch) Draw(buff draw.Image, xOff float64, yOff float64) {
	c.lock.RLock()
	c.subRenderables[c.curRenderable].Draw(buff, c.X()+xOff, c.Y()+yOff)
	c.lock.RUnlock()
}

// ShiftPos shifts the Switch's logical position
func (c *Switch) ShiftPos(x, y float64) {
	c.SetPos(c.X()+x, c.Y()+y)
}

// GetDims gets the current Renderables dimensions
func (c *Switch) GetDims() (int, int) {
	c.lock.RLock()
	w, h := c.subRenderables[c.curRenderable].GetDims()
	c.lock.RUnlock()
	return w, h
}

// Pause stops the current Renderable if possible
func (c *Switch) Pause() {
	c.lock.RLock()
	if cp, ok := c.subRenderables[c.curRenderable].(render.CanPause); ok {
		cp.Pause()
	}
	c.lock.RUnlock()
}

// Unpause tries to unpause the current Renderable if possible
func (c *Switch) Unpause() {
	c.lock.RLock()
	if cp, ok := c.subRenderables[c.curRenderable].(render.CanPause); ok {
		cp.Unpause()
	}
	c.lock.RUnlock()
}

// IsInterruptable returns whether the current renderable is interruptable
func (c *Switch) IsInterruptable() bool {
	c.lock.RLock()
	defer c.lock.RUnlock()
	if i, ok := c.subRenderables[c.curRenderable].(render.NonInterruptable); ok {
		return i.IsInterruptable()
	}
	return true
}

// IsStatic returns whether the current renderable is static
func (c *Switch) IsStatic() bool {
	c.lock.RLock()
	defer c.lock.RUnlock()
	if s, ok := c.subRenderables[c.curRenderable].(render.NonStatic); ok {
		return s.IsStatic()
	}
	return true
}

// SetTriggerID sets the ID AnimationEnd will trigger on for animating subtypes.
// Todo: standardize this with the other interface Set functions so that it
// also only acts on the current subRenderable, or the other way around, or
// somehow offer both options
func (c *Switch) SetTriggerID(cid event.CID) {
	c.lock.RLock()
	for _, r := range c.subRenderables {
		if t, ok := r.(render.Triggerable); ok {
			t.SetTriggerID(cid)
		}
	}
	c.lock.RUnlock()
}

// Revert will revert all parts of this Switch that can be reverted
func (c *Switch) Revert(mod int) {
	c.lock.RLock()
	for _, v := range c.subRenderables {
		switch t := v.(type) {
		case *render.Reverting:
			t.Revert(mod)
		}
	}
	c.lock.RUnlock()
}

// RevertAll will revert all parts of this Switch that can be reverted, back
// to their original state.
func (c *Switch) RevertAll() {
	c.lock.RLock()
	for _, v := range c.subRenderables {
		switch t := v.(type) {
		case *render.Reverting:
			t.RevertAll()
		}
	}
	c.lock.RUnlock()
}