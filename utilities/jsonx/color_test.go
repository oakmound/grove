package jsonx_test

import (
	"encoding/json"

	"bytes"
	"image/color"
	"testing"

	"github.com/oakmound/grove/utilities/jsonx"
)

func TestColor(t *testing.T) {
	type testCase struct {
		name          string
		input         color.RGBA
		expectedBytes []byte
		shouldErr     bool
	}
	tcs := []testCase{
		{
			name:          "basic",
			input:         color.RGBA{255, 0, 0, 255},
			expectedBytes: []byte("\"FF0000FF\""),
		},
	}
	for _, tc := range tcs {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			cr := jsonx.ColorRGBA(tc.input)
			outBytes, err := json.Marshal(cr)
			if err != nil {
				if !tc.shouldErr {
					t.Fatalf("got error when no error was expected: %v", err)
				}
			} else {
				if tc.shouldErr {
					t.Fatal("got no error when error was expected")
				}
				if !bytes.Equal(tc.expectedBytes, outBytes) {
					t.Fatalf("bytes expected %v vs got %v", string(tc.expectedBytes), string(outBytes))
				}
				c2 := &jsonx.ColorRGBA{}
				err := json.Unmarshal(outBytes, c2)
				if err != nil {
					// unmarshal after a successful marshal must not fail
					t.Fatalf("unmarshal failed: %v", err)
				}
				if *c2 != cr {
					t.Fatalf("unmarshal mismatched input")
				}
			}
		})
	}
}
