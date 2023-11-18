package image

import (
	"bytes"
	"image"
	"image/png"
	"testing"

	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/repository/image/mock"
	"github.com/go-park-mail-ru/2023_2_OND_team/pkg/logger"
	"github.com/go-park-mail-ru/2023_2_OND_team/pkg/validator/image/check"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestUploadInvalidImage(t *testing.T) {
	log, err := logger.New()
	if err != nil {
		t.Fatal(err)
	}
	usecase := New(log, nil)

	s, err := usecase.UploadImage("prefixPath", "image/png", 30, bytes.NewBuffer(nil), check.AnySize)
	require.Equal(t, ErrInvalidImage, err)
	require.Equal(t, "", s)
}

func TestUploadValidImage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	log, err := logger.New()
	if err != nil {
		t.Fatal(err)
	}

	imgRepo := mock.NewMockRepository(ctrl)
	usecase := New(log, imgRepo)
	rect := image.Rect(0, 0, 500, 500)
	img := image.NewRGBA(rect).SubImage(rect)
	buf := bytes.NewBuffer(nil)
	err = png.Encode(buf, img)
	require.NoError(t, err, "encode valid png image")

	prefixPath := "prefix"
	written := int64(buf.Len())
	expFilename := "image.png"
	imgRepo.EXPECT().
		SaveImage(prefixPath, "png", gomock.Any()).
		Return(expFilename, written, nil).
		Times(1)
	actualFile, err := usecase.UploadImage("prefix", "image/png", written, buf, check.AnySize)
	require.NoError(t, err)
	require.Equal(t, "https://pinspire.online:8081/"+expFilename, actualFile)
}
