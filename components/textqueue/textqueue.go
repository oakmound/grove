package textqueue

import (
	"fmt"
	"image/draw"
	"sync"
	"time"

	"github.com/oakmound/oak/v3/alg/floatgeom"
	"github.com/oakmound/oak/v3/dlog"
	"github.com/oakmound/oak/v3/entities/x/mods"
	"github.com/oakmound/oak/v3/event"
	"github.com/oakmound/oak/v3/render"
	"github.com/oakmound/oak/v3/render/mod"
	"github.com/oakmound/oak/v3/scene"
	"golang.org/x/image/colornames"
)

type queueItem struct {
	mod    render.Modifiable
	dropAt time.Time
}

// A TextQueue is a renderable entity that displays text in a column
// for a brief time before fading the text and dropping it. It accepts
// new text elements from the DisplayTextEvent event.
type TextQueue struct {
	event.CID
	render.LayeredPoint

	queueLock   sync.Mutex
	queue       []queueItem
	font        *render.Font
	sustainTime time.Duration
}

func (tq *TextQueue) Init() event.CID {
	return event.NextID(tq)
}

// New creates a customized TextQueue.
func New(ctx *scene.Context, initiatingEvent string, pos floatgeom.Point2, layer int, font *render.Font, sustainTime time.Duration) *TextQueue {
	tq := &TextQueue{}

	tq.CID = ctx.CallerMap.NextID(tq)

	tq.LayeredPoint = render.NewLayeredPoint(pos.X(), pos.Y(), layer)
	tq.font = font
	tq.queue = make([]queueItem, 0)
	tq.sustainTime = sustainTime

	if initiatingEvent == "" {
		initiatingEvent = DisplayTextEvent
	}

	tBind := func(id event.CID, payload interface{}) int {
		ent := ctx.CallerMap.GetEntity(id)
		tq, ok := ent.(*TextQueue)
		if !ok {
			dlog.Error("expected TextQueue, got " + fmt.Sprintf("%T", ent))
			return 1
		}
		str, ok := payload.(string)
		if !ok {
			dlog.Error("did not get string payload")
			return 0
		}
		r := tq.font.NewText(str, 0, 0)
		m := r.ToSprite().Modify(mods.HighlightOff(colornames.Black, 2, 1, 1))
		tq.queueLock.Lock()
		tq.queue = append([]queueItem{{
			mod:    m,
			dropAt: time.Now().Add(tq.sustainTime),
		}}, tq.queue...)
		tq.queueLock.Unlock()
		return 0
	}

	ctx.EventHandler.Bind(initiatingEvent, tq.CID, tBind)

	return tq
}

const DisplayTextEvent = "DisplayText"
const RecreationNeeded = "RecreationNeeded"

const yBuffer = 3

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

func (tq *TextQueue) GetDims() (int, int) {
	return 1, 1
}

func DisplayError(err error) {
	if err != nil {
		event.Trigger(DisplayTextEvent, err.Error())
	}
}
