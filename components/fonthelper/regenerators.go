package fonthelper

import (
	"image"
	"image/color"

	"github.com/oakmound/oak/v3/render"
)

// WithSize sets size on a font
func WithSize(size float64) func(render.FontGenerator) render.FontGenerator {
	return func(f render.FontGenerator) render.FontGenerator {
		f.Size = size
		return f
	}
}

// WithColor sets the color on a font
func WithColor(c color.Color) func(render.FontGenerator) render.FontGenerator {
	return func(f render.FontGenerator) render.FontGenerator {
		f.Color = image.NewUniform(c)
		return f
	}
}

// WithSizeAndColor sets the size and color on a font
func WithSizeAndColor(size float64, c color.Color) func(render.FontGenerator) render.FontGenerator {
	return func(f render.FontGenerator) render.FontGenerator {
		f.Color = image.NewUniform(c)
		f.Size = size
		return f
	}
}
