package textqueue

import (
	"fmt"
	"image/draw"
	"sync"
	"time"

	"github.com/oakmound/oak/v4/alg/floatgeom"
	"github.com/oakmound/oak/v4/dlog"

	"github.com/oakmound/oak/v4/event"
	"github.com/oakmound/oak/v4/render"
	"github.com/oakmound/oak/v4/render/mod"
	"github.com/oakmound/oak/v4/scene"
	"golang.org/x/image/colornames"
)

var (
	// TextQueuePublish: Triggered to publish an update to a specific text queue
	TextQueuePublish = event.RegisterEvent[string]()
)

type queueItem struct {
	mod    render.Modifiable
	dropAt time.Time
}

// A TextQueue is a renderable entity that displays text in a column
// for a brief time before fading the text and dropping it. It accepts
// new text elements from the DisplayTextEvent event.
type TextQueue struct {
	event.CallerID
	render.LayeredPoint

	queueLock   sync.Mutex
	queue       []queueItem
	font        *render.Font
	sustainTime time.Duration
}

func (tq *TextQueue) CID() event.CallerID {
	return tq.CallerID.CID()
}

// New creates a customized TextQueue.
func New(ctx *scene.Context, registeredEvents []event.UnsafeEventID, pos floatgeom.Point2, layer int, font *render.Font, sustainTime time.Duration) *TextQueue {
	tq := &TextQueue{}
	tq.CallerID = ctx.Register(tq)

	tq.LayeredPoint = render.NewLayeredPoint(pos.X(), pos.Y(), layer)
	tq.font = font
	tq.queue = make([]queueItem, 0)
	tq.sustainTime = sustainTime

	event.Bind(ctx, TextQueuePublish, tq, PrintBind)

	if len(registeredEvents) == 0 {
		return tq
	}

	unsafePrintBind := func(id event.CallerID, handler event.Handler, payload interface{}) event.Response {
		ent := handler.GetCallerMap().GetEntity(id)
		_, ok := ent.(*TextQueue)
		if !ok && ent != nil {
			dlog.Error("expected TextQueue, got " + fmt.Sprintf("%T", ent))
			return 1
		}
		PrintBind(tq, fmt.Sprintf("%v", payload))

		return 0
	}

	for _, re := range registeredEvents {
		ctx.UnsafeBind(re, tq.CID(), unsafePrintBind)
	}

	return tq
}

func PrintBind(tq *TextQueue, str string) event.Response {
	r := tq.font.NewText(str, 0, 0)
	m := r.ToSprite().Modify(mod.HighlightOff(colornames.Black, 2, 1, 1))
	tq.queueLock.Lock()
	tq.queue = append([]queueItem{{
		mod:    m,
		dropAt: time.Now().Add(tq.sustainTime),
	}}, tq.queue...)
	tq.queueLock.Unlock()
	return 0
}

const DisplayTextEvent = "DisplayText"

const yBuffer = 3

// Draw the textqueue's contents
func (tq *TextQueue) Draw(buff draw.Image, xOff, yOff float64) {
	if len(tq.queue) == 0 {
		return
	}
	now := time.Now()
	secondFromNow := now.Add(1 * time.Second)

	if now.After(tq.queue[len(tq.queue)-1].dropAt) {
		tq.queueLock.Lock()
		tq.queue = tq.queue[:len(tq.queue)-1]
		tq.queueLock.Unlock()
	}

	xOff += tq.X()
	yOff += tq.Y()
	for _, item := range tq.queue {
		_, y := item.mod.GetDims()
		item.mod.Draw(buff, xOff, yOff)
		if secondFromNow.After(item.dropAt) {
			item.mod.Filter(mod.Fade(5))
		}
		yOff += float64(y) + yBuffer
	}
}

// GetDims needs to have some size so give it the minimal one.
func (tq *TextQueue) GetDims() (int, int) {
	return 1, 1
}
