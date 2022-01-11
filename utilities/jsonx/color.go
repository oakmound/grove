package jsonx

import (
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"image/color"
	"strings"

	"golang.org/x/image/colornames"
)

const rgbaHexFormat = "%02X%02X%02X%02X"

type ColorRGBA color.RGBA

func (c ColorRGBA) MarshalJSON() ([]byte, error) {
	s := fmt.Sprintf(rgbaHexFormat, c.R, c.G, c.B, c.A)
	return json.Marshal(s)
}

func (c *ColorRGBA) UnmarshalJSON(b []byte) error {
	var v interface{}
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	switch value := v.(type) {
	case string:
		s := strings.ToLower(value)
		if c2, ok := colornames.Map[s]; ok {
			*c = ColorRGBA(c2)
		} else {
			n, err := fmt.Sscanf(value, rgbaHexFormat, &c.R, &c.G, &c.B, &c.A)
			if n != 4 {
				return fmt.Errorf("not enough hex values (%v)", n)
			}
			if err != nil {
				return fmt.Errorf("hex scan failed: %w", err)
			}
		}
		return nil
	default:
		return errors.New("invalid color name")
	}
}

type ColorUniform image.Uniform

func (c ColorUniform) MarshalJSON() ([]byte, error) {
	r, g, b, a := c.C.RGBA()
	s := fmt.Sprintf(rgbaHexFormat, r/255, g/255, b/255, a/255)
	return []byte(s), nil
}

func (c *ColorUniform) UnmarshalJSON(b []byte) error {
	var v interface{}
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	switch value := v.(type) {
	case string:
		rgba := color.RGBA{}
		s := strings.ToLower(value)
		if c2, ok := colornames.Map[s]; ok {
			rgba = c2
		} else {
			n, err := fmt.Sscanf(value, rgbaHexFormat, &rgba.R, &rgba.G, &rgba.B, &rgba.A)
			if n != 4 {
				return fmt.Errorf("not enough hex values (%v)", n)
			}
			if err != nil {
				return fmt.Errorf("hex scan failed: %w", err)
			}
		}
		*c = ColorUniform(*image.NewUniform(rgba))
		return nil
	default:
		return errors.New("invalid color name")
	}
}
