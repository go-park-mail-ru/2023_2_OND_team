package image

import (
	"fmt"
	"io"
	"strings"

	"github.com/go-park-mail-ru/2023_2_OND_team/pkg/validator/image/jpeg"
	"github.com/go-park-mail-ru/2023_2_OND_team/pkg/validator/image/png"
	"github.com/go-park-mail-ru/2023_2_OND_team/pkg/validator/image/svg"
	"github.com/go-park-mail-ru/2023_2_OND_team/pkg/validator/image/webp"
)

func IsValidImage(image io.Reader, mimeType string) (string, bool) {
	if !strings.HasPrefix(mimeType, "image/") {
		return "", false
	}
	mimeType = mimeType[6:]

	switch {
	case strings.HasPrefix(mimeType, "jpeg"):
		return "jpg", jpeg.IsValidJPEG(image)
	case strings.HasPrefix(mimeType, "png"):
		return "png", png.IsValidPNG(image)
	case strings.HasPrefix(mimeType, "svg"):
		fmt.Println("SVG")
		return "svg", svg.IsValidSVG(image)
	case strings.HasPrefix(mimeType, "webp"):
		return "webp", webp.IsValidWEBP(image)
	default:
		return "", false
	}

}
