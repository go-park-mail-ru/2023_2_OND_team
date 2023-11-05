package svg

import (
	"io"

	"github.com/go-park-mail-ru/2023_2_OND_team/pkg/validator/image/check"
	"github.com/tdewolff/canvas"
)

func IsValidSVG(r io.Reader, check check.CheckSize) bool {
	can, err := canvas.ParseSVG(r)
	if err != nil {
		return false
	}

	return check(can.W, can.H)
}
