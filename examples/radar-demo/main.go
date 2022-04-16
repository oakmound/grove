package main

import (
	"image/color"
	"math"
	"math/rand"

	"github.com/oakmound/oak/v3/alg/floatgeom"
	"github.com/oakmound/oak/v3/alg/intgeom"

	oak "github.com/oakmound/oak/v3"
	"github.com/oakmound/oak/v3/entities"
	"github.com/oakmound/oak/v3/entities/x/move"
	"github.com/oakmound/oak/v3/event"

	"github.com/oakmound/oak/v3/physics"
	"github.com/oakmound/oak/v3/render"
	"github.com/oakmound/oak/v3/scene"

	"github.com/oakmound/grove/components/radar"
)

const (
	xLimit = 1000
	yLimit = 1000
)

// This example demonstrates making a basic radar or other custom renderable
// type. The radar here acts as a UI element, staying on screen, and follows
// around a player character.
//TODO: Remove and or link to grove radar as it is cleaner
// https://github.com/oakmound/grove/tree/master/components/radar

func main() {
	oak.AddScene("demo", scene.Scene{Start: func(ctx *scene.Context) {
		render.Draw(render.NewDrawFPS(0.03, nil, 10, 10))

		char := entities.NewMoving(200, 200, 50, 50, render.NewColorBox(50, 50, color.RGBA{125, 125, 0, 255}), nil, 0, 1)
		char.Speed = physics.NewVector(3, 3)

		oak.SetViewportBounds(intgeom.NewRect2(0, 0, xLimit, yLimit))
		moveRect := floatgeom.NewRect2(0, 0, xLimit, yLimit)
		event.Bind(ctx, event.Enter, char, func(char *entities.Moving, ev event.EnterPayload) event.Response {
			move.WASD(char)
			move.Limit(char, moveRect)
			move.CenterScreenOn(char)
			return 0
		})
		render.Draw(char.R, 1, 2)

		// Create the Radar
		center := radar.Point{X: char.Xp(), Y: char.Yp()}
		points := make(map[radar.Point]color.Color)
		w := 100
		h := 100
		r := radar.NewRadar(w, h, points, center, 10)
		r.SetPos(float64(ctx.Window.Width()-w), 0)

		for i := 0; i < 5; i++ {
			x, y := rand.Float64()*400, rand.Float64()*400
			enemy := newEnemyOnRadar(ctx, x, y)
			event.Bind(ctx, event.Enter, enemy, standardEnemyMove)
			render.Draw(enemy.R, 1, 1)
			r.AddPoint(radar.Point{X: enemy.Xp(), Y: enemy.Yp()}, color.RGBA{255, 255, 0, 255})
		}

		render.Draw(r, 2)

		for x := 0; x < xLimit; x += 64 {
			for y := 0; y < yLimit; y += 64 {
				r := uint8(rand.Intn(120))
				b := uint8(rand.Intn(120))
				cb := render.NewColorBox(64, 64, color.RGBA{r, 0, b, 255})
				cb.SetPos(float64(x), float64(y))
				render.Draw(cb, 0)
			}
		}

	}})

	render.SetDrawStack(
		render.NewCompositeR(),
		render.NewDynamicHeap(),
		render.NewStaticHeap(),
	)
	oak.Init("demo")
}

type enemyOnRadar struct {
	*entities.Moving
}

func newEnemyOnRadar(ctx *scene.Context, x, y float64) *enemyOnRadar {
	eor := new(enemyOnRadar)
	eor.Moving = entities.NewMoving(50, y, 50, 50, render.NewColorBox(25, 25, color.RGBA{0, 200, 0, 0}), nil, ctx.Register(eor), 0)
	eor.Speed = physics.NewVector(-1*(rand.Float64()*2+1), rand.Float64()*2-1)
	eor.Delta = eor.Speed
	return eor
}

func standardEnemyMove(eor *enemyOnRadar, ev event.EnterPayload) event.Response {
	if eor.X() < 0 {
		eor.Delta.SetPos(math.Abs(eor.Speed.X()), (eor.Speed.Y()))
	}
	if eor.X() > xLimit-eor.W {
		eor.Delta.SetPos(-1*math.Abs(eor.Speed.X()), (eor.Speed.Y()))
	}
	if eor.Y() < 0 {
		eor.Delta.SetPos(eor.Speed.X(), math.Abs(eor.Speed.Y()))
	}
	if eor.Y() > yLimit-eor.H {
		eor.Delta.SetPos(eor.Speed.X(), -1*math.Abs(eor.Speed.Y()))
	}
	eor.ShiftX(eor.Delta.X())
	eor.ShiftY(eor.Delta.Y())
	return 0
}
