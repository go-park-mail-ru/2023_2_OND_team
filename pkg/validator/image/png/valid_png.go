package png

import (
	"image/png"
	"io"
)

func IsValidPNG(r io.Reader) bool {
	cfg, err := png.DecodeConfig(r)
	if err != nil {
		return false
	}

	if cfg.Height < 200 || cfg.Width < 200 || cfg.Height > 600 || cfg.Width > 600 {
		return false
	}
	return true
}
