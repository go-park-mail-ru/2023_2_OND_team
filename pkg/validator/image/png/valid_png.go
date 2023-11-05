package png

import (
	"image/png"
	"io"

	"github.com/go-park-mail-ru/2023_2_OND_team/pkg/validator/image/check"
)

func IsValidPNG(r io.Reader, check check.CheckSize) bool {
	cfg, err := png.DecodeConfig(r)
	if err != nil {
		return false
	}

	return check(float64(cfg.Width), float64(cfg.Height))
}
