package textfit

import (
	"testing"

	"github.com/oakmound/oak/v3/alg/floatgeom"
	"github.com/oakmound/oak/v3/render"
)

func TestNew(t *testing.T) {
	type testCase struct {
		opts []Option
		// should not error
	}
	tcs := []testCase{
		{
			opts: []Option{
				String("The final episode of \"The Legend of High School\" show has leaked two hours before its premiere. As a moderator of the official \"The Legend of High School\" forum, you need to <b>keep the fans from learning any details about the ending of the show</b>.You've read the books before they made it a show, so you're already familiar with the plot points that are going to happen: <b>Alexzandre will Confess to Jeremiah</b>, <b>Olivette will miss the Polar Dance and not be sad about it</b>, <b> Tim Runnings will Tie with Ran Jennings for first place in the Quinqometry exams </b>, and finally <b>Sam Sam will show up to Graduation, Walk and end the series with his Goodbye Speech</b>. As always, this is a family friendly forum. <b>Discussion of non-canon romantic pairings of characters on the show is not allowed</b>. <b>Profanity of -any- kind is also not allowed</b>.\n<b>Linking to the leaked episode is not allowed </b>. <b> Asking for links to the leaked episode is also not allowed </b>."),
				MinSize(1),
				MaxSize(22),
				Font(render.DefaultFont()),
				Dimensions(floatgeom.Point2{300, 300}),
			},
		},
	}
	for _, tc := range tcs {
		_, err := New(tc.opts...)
		if err != nil {
			t.Fatalf("got error: %v", err)
		}
	}
}
