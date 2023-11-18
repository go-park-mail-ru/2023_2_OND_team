package jpeg

import (
	"image/jpeg"
	"io"

	"github.com/go-park-mail-ru/2023_2_OND_team/pkg/validator/image/check"
)

func IsValidJPEG(r io.Reader, check check.CheckSize) bool {
	cfg, err := jpeg.DecodeConfig(r)
	if err != nil {
		return false
	}
	return check(float64(cfg.Width), float64(cfg.Height))
}
