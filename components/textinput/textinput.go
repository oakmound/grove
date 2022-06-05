package textinput

import (
	"fmt"
	"image/color"
	"sync"
	"time"

	"github.com/oakmound/oak/v4/alg/floatgeom"
	"github.com/oakmound/oak/v4/entities"
	"github.com/oakmound/oak/v4/event"
	"github.com/oakmound/oak/v4/key"
	"github.com/oakmound/oak/v4/mouse"
	"github.com/oakmound/oak/v4/render"
	"github.com/oakmound/oak/v4/scene"
)

// TextInput provides a nicer way to handle input of text
// Notably it creates a blinking input cursor
type TextInput struct {
	id event.CallerID
	*entities.Entity
	ctx *scene.Context

	bindingLock sync.Mutex

	textLock    sync.Mutex
	currentText *string

	editing bool
	w, h    float64

	position   floatgeom.Point2
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

	entityOptions []entities.Option
}

func (ti *TextInput) CID() event.CallerID {
	return ti.id
}

// New textinput for the scene given a set of options
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
	ti.id = ctx.Register(ti)
	for _, opt := range opts {
		opt(ti)
	}
	ti.font = ti.font.Copy()
	r := ti.font.NewStrPtrText(ti.currentText, 0, 0)
	// ti.Solid = entities.NewSolid(ti.x, ti.y, ti.w, ti.h, r, mouse.DefaultTree, ti.Init())
	// ti.Entity =

	ti.Entity = entities.New(ctx,
		append(ti.entityOptions,
			entities.WithParent(ti),
			entities.WithRenderable(r),
			entities.WithDimensions(floatgeom.Point2{ti.w, ti.h}),
			entities.WithPosition(ti.position),
			entities.WithUseMouseTree(true))...,
	)
	event.Bind(ctx, mouse.PressOn, ti, func(textIn *TextInput, me *mouse.Event) event.Response {

		return ti.startTyping(ctx, me)
	})

	return ti
}

// startTyping bind initiates the ability to add text to the textinput area
func (ti *TextInput) startTyping(ctx *scene.Context, me *mouse.Event) event.Response {

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
	// ti.updateBlinkerToMouse(me)
	editB1 := event.Bind(ctx, key.AnyDown, ti, editBinding)
	editB2 := event.Bind(ctx, key.AnyHeld, ti, editBinding)
	fmt.Println("starting to type")
	// for some reason unbind isnt working as expected below so paranoid unbind
	var stopB1 event.Binding
	stopB1 = event.Bind(ctx, mouse.Release, ti, func(textIn *TextInput, me *mouse.Event) event.Response {
		fmt.Println("Stopping your typing")
		// textIn.stopTyping(ctx, editB1, editB2)
		// ctx.Unbind(stopB1)
		return event.ResponseUnbindThisBinding
	})
	if editB1 == editB2 && stopB1 == editB1 {
		fmt.Println("Just having these things being used and yes this should never happen ")
	}
	// ti.CheckedBind(mouse.Click, func(ti *TextInput, ev interface{}) int {
	// 	return ti.stopTyping()
	// })

	return event.ResponseUnbindThisBinding
}

func (ti *TextInput) stopTyping(ctx *scene.Context, b1, b2 event.Binding) event.Response {
	ti.bindingLock.Lock()
	defer ti.bindingLock.Unlock()
	// only stop editing if not already editing
	if !ti.editing {
		fmt.Println("hi")
		return event.ResponseUnbindThisBinding
	}

	ti.editing = false
	// ti.undrawBlinker()
	if ti.finalizer != nil {
		if ti.sensitive {
			ti.finalizer(ti.sensitiveText)
		} else {
			ti.finalizer(*ti.currentText)
		}
	}
	// event.Bind(ctx, mouse.PressOn, ti, func(textIn *TextInput, me *mouse.Event) event.Response {
	// 	return textIn.startTyping(ctx, me)
	// })

	ctx.Unbind(b1)
	ctx.Unbind(b2)

	return event.ResponseUnbindThisBinding
}

func editBinding(textIn *TextInput, ke key.Event) event.Response {

	// safety check that we are actually editing
	if !textIn.editing {
		return event.ResponseUnbindThisBinding
	}

	textIn.textLock.Lock()
	txt := *textIn.currentText
	textIn.textLock.Unlock()

	code := ke.Code.String()[4:]
	shift := 0

	switch code {
	default:
		if textIn.sensitive {
			txt += "*"
			textIn.sensitiveText = textIn.sensitiveText[:textIn.blinkerIndex] + string(ke.Rune) + textIn.sensitiveText[textIn.blinkerIndex:]
		} else {
			txt = txt[:textIn.blinkerIndex] + string(ke.Rune) + txt[textIn.blinkerIndex:]
		}
		shift = len(string(ke.Rune))
	}

	textIn.textLock.Lock()
	*textIn.currentText = txt
	textIn.textLock.Unlock()
	fmt.Println(shift)
	// textIn.updateBlinkerRelative(shift)

	return 0
}
