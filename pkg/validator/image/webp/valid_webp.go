package webp

import (
	"io"

	"github.com/go-park-mail-ru/2023_2_OND_team/pkg/validator/image/check"
	"golang.org/x/image/webp"
)

func IsValidWEBP(r io.Reader, check check.CheckSize) bool {
	cfg, err := webp.DecodeConfig(r)
	if err != nil {
		return false
	}
	return check(float64(cfg.Width), float64(cfg.Height))
}
