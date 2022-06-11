package textfit

import (
	"fmt"
	"sort"
	"strings"

	"github.com/oakmound/oak/v4/alg/floatgeom"
	"github.com/oakmound/oak/v4/render"
)

type Generator struct {
	Text       string
	Font       *render.Font
	Dimensions floatgeom.Point2
	MinSize    int
	MaxSize    int
	BreakStyle
	RelativePos floatgeom.Point2
	// overflow style - clip / newlines on word / newlines on character / trail off (...)
	// padding
}

type BreakStyle byte

const (
	BreakStyleCharacter BreakStyle = iota
	BreakStyleWord
)

func (g *Generator) generate() (render.Modifiable, error) {
	size := g.MaxSize
	dims := g.Dimensions.Sub(g.RelativePos.MulConst(2))

	var sections []string
ONESECTION:
	for {
		font, _ := g.Font.RegenerateWith(func(f render.FontGenerator) render.FontGenerator {
			f.Size = float64(size)
			return f
		})
		sections = []string{g.Text}
	HEIGHTCHECK:
		for {
			remainingHeight := dims.Y() - float64(size*len(sections))
			if remainingHeight < 0 {
				size--
				if size < g.MinSize {
					return nil, fmt.Errorf("size fell below minSize; decrease min size, decrease text ct, or increase bounds")
				}
				continue ONESECTION
			}
			for i, sec := range sections {
				sec = strings.TrimSpace(sec)
				length := font.MeasureString(sec).Round()
				if length > int(dims.X()) {
					breakPoint := sort.Search(len(sec), func(i int) bool {
						// returns the smallest index for which i is true,
						// so we want to return true if we are above our goal width,
						// and then subtract 1.
						length := font.MeasureString(sec[:i]).Round()
						return length > int(dims.X())
					})
					breakPoint--
					if g.BreakStyle == BreakStyleWord {
						// walk back to the last space
						for breakPoint > 1 && sec[breakPoint] != ' ' {
							breakPoint--
						}
					}

					sections[i] = sec[:breakPoint]
					sections = append(sections[:i+1], append(
						[]string{sec[breakPoint:]}, sections[i+1:]...)...)

					continue HEIGHTCHECK
				}
			}
			break ONESECTION
		}
	}

	comp := render.NewCompositeR()

	font, _ := g.Font.RegenerateWith(func(f render.FontGenerator) render.FontGenerator {
		f.Size = float64(size)
		return f
	})
	x := g.RelativePos.X()
	y := g.RelativePos.Y()

	for _, sec := range sections {
		sec = strings.TrimSpace(sec)
		t := font.NewText(sec, x, y)
		y += float64(size)
		comp.Append(t.ToSprite())
	}

	return comp.ToSprite(), nil
}

func defaultGenerator() *Generator {
	return &Generator{
		MinSize:    5,
		MaxSize:    30,
		Font:       render.DefaultFont(),
		Text:       "placeholder",
		Dimensions: floatgeom.Point2{50, 50},
	}
}

func MustNew(options ...Option) render.Modifiable {
	r, err := New(options...)
	if err != nil {
		panic(err)
	}
	return r
}

func New(options ...Option) (render.Modifiable, error) {
	gen := defaultGenerator()
	for _, opt := range options {
		opt(gen)
	}
	return gen.generate()
}

type Option (func(*Generator))

func String(s string) Option {
	return func(g *Generator) {
		g.Text = s
	}
}

func Font(f *render.Font) Option {
	return func(g *Generator) {
		g.Font = f
	}
}

func Dimensions(p floatgeom.Point2) Option {
	return func(g *Generator) {
		g.Dimensions = p
	}
}

func Inset(w, h float64) Option {
	// reduce dimensions by w*2 and h*2
	// set relpos to w,h
	return func(g *Generator) {
		g.RelativePos = floatgeom.Point2{w, h}
	}
}

func MinSize(min int) Option {
	return func(g *Generator) {
		g.MinSize = min
	}
}

func MaxSize(max int) Option {
	return func(g *Generator) {
		g.MaxSize = max
	}
}

func WithBreakStyle(bs BreakStyle) Option {
	return func(g *Generator) {
		g.BreakStyle = bs
	}
}
