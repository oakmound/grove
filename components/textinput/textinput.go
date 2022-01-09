package textinput

import (
	"fmt"
	"image/color"
	"sync"
	"time"

	"github.com/oakmound/oak/v3/alg/floatgeom"
	"github.com/oakmound/oak/v3/dlog"
	"github.com/oakmound/oak/v3/entities"
	"github.com/oakmound/oak/v3/event"
	"github.com/oakmound/oak/v3/key"
	"github.com/oakmound/oak/v3/mouse"
	"github.com/oakmound/oak/v3/render"
	"github.com/oakmound/oak/v3/scene"
	"github.com/oakmound/oak/v3/timing"
)

type TextInput struct {
	*entities.Solid
	ctx *scene.Context

	bindingLock sync.Mutex

	textLock    sync.Mutex
	currentText *string

	editing    bool
	x, y       float64
	w, h       float64
	textOffset floatgeom.Point2

	finalizer   func(string)
	onFirstEdit func(ti *TextInput)
	onEdit      func(ti *TextInput)

	font *render.Font

	blinkerLock   sync.Mutex
	blinker       render.Renderable
	blinkerColor  color.Color
	blinkRate     time.Duration
	blinkerIndex  int
	blinkerLayers []int

	sensitive     bool
	sensitiveText string
}

func New(ctx *scene.Context, opts ...Option) *TextInput {
	emptyString := ""
	ti := &TextInput{
		ctx:           ctx,
		w:             100,
		h:             20,
		font:          render.DefaultFont(),
		blinkerColor:  color.RGBA{255, 255, 255, 255},
		currentText:   &emptyString,
		blinkerLayers: []int{0, 2},
	}
	for _, opt := range opts {
		opt(ti)
	}
	ti.font = ti.font.Copy()
	r := ti.font.NewStrPtrText(ti.currentText, 0, 0)
	ti.Solid = entities.NewSolid(ti.x, ti.y, ti.w, ti.h, r, mouse.DefaultTree, ti.Init())
	ti.bindStartTyping()
	ti.R.SetPos(ti.x+ti.textOffset.X(), ti.y+ti.textOffset.Y())
	return ti
}

// Select and Deselect simulate mouse click actions to enable or disable
// typing in a text input.

func (ti *TextInput) Select() {
	ti.Trigger(mouse.ClickOn, mouse.Event{})
}

func (ti *TextInput) Deselect() {
	ti.stopTyping()
}

func (ti *TextInput) startTyping(me mouse.Event) int {
	ti.bindingLock.Lock()
	defer ti.bindingLock.Unlock()

	if ti.onFirstEdit != nil {
		ti.onFirstEdit(ti)
		ti.onFirstEdit = nil
	}
	if ti.onEdit != nil {
		ti.onEdit(ti)
	}
	ti.editing = true
	ti.updateBlinkerToMouse(me)
	ti.Bind(key.Down, editBinding)
	ti.Bind(key.Held, editBinding)
	ti.CheckedBind(mouse.Click, func(ti *TextInput, ev interface{}) int {
		return ti.stopTyping()
	})

	return event.UnbindSingle
}

func (ti *TextInput) bindStartTyping() {
	ti.CheckedBind(mouse.ClickOn, func(ti *TextInput, ev interface{}) int {
		me, ok := ev.(*mouse.Event)
		if !ok {
			fmt.Println("text input received non-mouse event argument to ClickOn")
			return 0
		}
		return ti.startTyping(*me)
	})
}

func (ti *TextInput) stopTyping() int {
	ti.bindingLock.Lock()
	defer ti.bindingLock.Unlock()

	if ti.editing {
		ti.editing = false
		ti.undrawBlinker()
		if ti.finalizer != nil {
			if ti.sensitive {
				ti.finalizer(ti.sensitiveText)
			} else {
				ti.finalizer(*ti.currentText)
			}
		}
		ti.bindStartTyping()
		event.UnbindBindable(event.UnbindOption{
			Event: event.Event{
				Name:     key.Down,
				CallerID: ti.CID,
			},
			Fn: editBinding,
		})
		event.UnbindBindable(event.UnbindOption{
			Event: event.Event{
				Name:     key.Held,
				CallerID: ti.CID,
			},
			Fn: editBinding,
		})
		return event.UnbindSingle
	}
	return 0
}

func (ti *TextInput) Init() event.CID {
	return event.NextID(ti)
}

func (ti *TextInput) CheckedBind(name string, f func(*TextInput, interface{}) int) {
	ti.Bind(name, func(id event.CID, ev interface{}) int {
		ti, ok := id.E().(*TextInput)
		if !ok {
			dlog.Error("Non-TextInput passed to TextInput binding")
			return 0
		}
		return f(ti, ev)
	})
}

func (ti *TextInput) updateBlinkerToMouse(me mouse.Event) {
	ti.textLock.Lock()
	// convert me to index position
	// linear scan until its demonstrated we need something with better performance
	var textIndex int
	for i := 0; i < len(*ti.currentText); i++ {
		charX := float64(ti.font.MeasureString((*ti.currentText)[:i]).Round())
		charX += ti.R.X()
		if charX > me.X() {
			textIndex = i
			break
		}
	}
	ti.textLock.Unlock()

	ti.updateBlinker(textIndex)
}

func (ti *TextInput) updateBlinkerRelative(shift int) {
	ti.updateBlinker(ti.blinkerIndex + shift)
}

func (ti *TextInput) updateBlinker(textIndex int) {
	ti.blinkerLock.Lock()
	defer ti.blinkerLock.Unlock()
	if ti.blinker != nil {
		ti.blinker.Undraw()
	}
	ti.textLock.Lock()
	var w float64
	h := ti.font.Height()
	if textIndex < 0 {
		w = 0
		ti.blinkerIndex = 0
	} else {
		if textIndex >= len(*ti.currentText) {
			textIndex = len(*ti.currentText)
		}
		fixedWidth := ti.font.MeasureString((*ti.currentText)[:textIndex])
		w = float64(fixedWidth.Round())
		ti.blinkerIndex = textIndex
	}
	x, y := ti.R.X(), ti.R.Y()
	ti.textLock.Unlock()
	if ti.blinkRate != 0 {
		ti.blinker = render.NewSequence(timing.FrameDelayToFPS(ti.blinkRate),
			render.NewLine(x+w, y, x+w, y+h, ti.blinkerColor),
			render.EmptyRenderable(),
		)
	} else {
		ti.blinker = render.NewLine(x+w, y, x+w, y+h, ti.blinkerColor)
	}
	ti.ctx.DrawStack.Draw(ti.blinker, ti.blinkerLayers...)
}

func (ti *TextInput) undrawBlinker() {
	ti.blinkerLock.Lock()
	defer ti.blinkerLock.Unlock()
	if ti.blinker != nil {
		ti.blinker.Undraw()
	}
}

// TODO: checked bindings can't be unbound directly

func editBinding(id event.CID, ev interface{}) int {
	ti, ok := id.E().(*TextInput)
	if !ok {
		return event.UnbindSingle
	}
	if !ti.editing {
		return event.UnbindSingle
	}
	k, ok := ev.(key.Event)
	if !ok {
		dlog.Error("Got non key event in text input edit")
		return 0
	}
	ti.textLock.Lock()
	txt := *ti.currentText
	ti.textLock.Unlock()

	code := k.Code.String()[4:]
	shift := 0
	switch code {
	case key.Enter, key.Escape:
		ti.bindingLock.Lock()
		defer ti.bindingLock.Unlock()
		ti.editing = false
		ti.undrawBlinker()
		if ti.finalizer != nil {
			if ti.sensitive {
				ti.finalizer(ti.sensitiveText)
			} else {
				ti.finalizer(txt)
			}
		}
		ti.bindStartTyping()
		return event.UnbindSingle
	case key.DeleteBackspace:
		if len(txt) != 0 && ti.blinkerIndex != 0 {
			if ti.blinkerIndex >= len(txt) {
				txt = txt[:ti.blinkerIndex-1]
			} else {
				txt = txt[:ti.blinkerIndex-1] + txt[ti.blinkerIndex:]
			}
		}
		if ti.sensitive && len(ti.sensitiveText) != 0 && ti.blinkerIndex != 0 {
			if ti.blinkerIndex >= len(ti.sensitiveText) {
				ti.sensitiveText = ti.sensitiveText[:ti.blinkerIndex-1]
			} else {
				ti.sensitiveText = ti.sensitiveText[:ti.blinkerIndex-1] + ti.sensitiveText[ti.blinkerIndex:]
			}
		}
		shift = -1
	case key.LeftShift, key.RightShift, key.Tab:
	case key.LeftArrow:
		ti.updateBlinkerRelative(-1)
		return 0
	case key.RightArrow:
		ti.updateBlinkerRelative(1)
		return 0
	default:
		if ti.sensitive {
			txt += "*"
			ti.sensitiveText = ti.sensitiveText[:ti.blinkerIndex] + string(k.Rune) + ti.sensitiveText[ti.blinkerIndex:]
		} else {
			txt = txt[:ti.blinkerIndex] + string(k.Rune) + txt[ti.blinkerIndex:]
		}
		shift = len(string(k.Rune))
	}
	ti.textLock.Lock()
	*ti.currentText = txt
	ti.textLock.Unlock()
	ti.updateBlinkerRelative(shift)
	return 0
}
