package image

import (
	"fmt"
	"io"
	"strings"

	"github.com/go-park-mail-ru/2023_2_OND_team/pkg/validator/image/check"
	"github.com/go-park-mail-ru/2023_2_OND_team/pkg/validator/image/jpeg"
	"github.com/go-park-mail-ru/2023_2_OND_team/pkg/validator/image/png"
	"github.com/go-park-mail-ru/2023_2_OND_team/pkg/validator/image/svg"
	"github.com/go-park-mail-ru/2023_2_OND_team/pkg/validator/image/webp"
)

func IsValidImage(image io.Reader, mimeType string, check check.CheckSize) (string, bool) {
	if !strings.HasPrefix(mimeType, "image/") {
		return "", false
	}
	mimeType = mimeType[6:]

	switch {
	case strings.HasPrefix(mimeType, "jpeg"):
		return "jpg", jpeg.IsValidJPEG(image, check)
	case strings.HasPrefix(mimeType, "png"):
		return "png", png.IsValidPNG(image, check)
	case strings.HasPrefix(mimeType, "svg"):
		fmt.Println("SVG")
		return "svg", svg.IsValidSVG(image, check)
	case strings.HasPrefix(mimeType, "webp"):
		return "webp", webp.IsValidWEBP(image, check)
	default:
		return "", false
	}

}
