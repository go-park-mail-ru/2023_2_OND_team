package svg

import (
	"io"

	"github.com/tdewolff/canvas"
)

func IsValidSVG(r io.Reader) bool {
	can, err := canvas.ParseSVG(r)
	if err != nil {
		return false
	}

	if can.H < 200 || can.W < 200 || can.H > 600 || can.W > 600 {
		return false
	}

	return true
}
