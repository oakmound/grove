package textinput

import (
	"image/color"
	"time"

	"github.com/oakmound/oak/v4/alg/floatgeom"
	"github.com/oakmound/oak/v4/render"
)

// Option for configuring a TextInput
type Option func(*TextInput)

// And together the options
func And(opts ...Option) Option {
	return func(t *TextInput) {
		for _, opt := range opts {
			opt(t)
		}
	}
}

// WithStr sets the str as the text on the textinput
func WithStr(s string) Option {
	return func(t *TextInput) {
		t.currentText = &s
	}
}

func WithStrPtr(s *string) Option {
	return func(t *TextInput) {
		t.currentText = s
	}
}

func WithPlaceholder(s string) Option {
	return func(t *TextInput) {
		t.currentText = &s
		t.onFirstEdit = func(ti *TextInput) {
			*t.currentText = ""
		}
	}
}

func WithPos(x, y float64) Option {
	return func(t *TextInput) {
		t.x = x
		t.y = y
	}
}

func WithDims(w, h float64) Option {
	return func(t *TextInput) {
		t.w = w
		t.h = h
	}
}

func WithFinalizer(f func(string)) Option {
	return func(t *TextInput) {
		t.finalizer = f
	}
}

func WithFont(f *render.Font) Option {
	return func(t *TextInput) {
		t.font = f
	}
}

func WithBlinkerColor(c color.Color) Option {
	return func(t *TextInput) {
		t.blinkerColor = c
	}
}

func WithSensitive(sensitive bool) Option {
	return func(t *TextInput) {
		t.sensitive = sensitive
	}
}

func WithBlinkRate(blinkRate time.Duration) Option {
	return func(t *TextInput) {
		t.blinkRate = blinkRate
	}
}

func WithTextOffset(x, y float64) Option {
	return func(t *TextInput) {
		t.textOffset = floatgeom.Point2{x, y}
	}
}

func WithOnEdit(onEdit func(*TextInput)) Option {
	return func(t *TextInput) {
		t.onEdit = onEdit
	}
}

func WithBlinkerLayers(layers ...int) Option {
	return func(t *TextInput) {
		t.blinkerLayers = layers
	}
}

// func WithEntityOptions(opts ...entities.Option) Option {
// 	return func(t *TextInput) {
// 		t.entityOptions = opts
// 	}
// }
