package image

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"github.com/google/uuid"
)

//go:generate mockgen -destination=./mock/image_mock.go -package=mock -source=repo.go Repository
type Repository interface {
	SaveImage(prefixPath, extension string, image io.Reader) (filename string, written int64, err error)
	SetBasePath(path string)
	SetDirToSave(fn func() string)
}

type imageRepoFS struct {
	basePath        string
	m               sync.Mutex
	directoryToSave func() string
}

func NewImageRepoFS(basePath string) *imageRepoFS {
	return &imageRepoFS{
		basePath:        basePath,
		m:               sync.Mutex{},
		directoryToSave: dirToSave,
	}
}

func (img *imageRepoFS) SaveImage(prefixPath, extension string, image io.Reader) (filename string, written int64, err error) {
	id, err := uuid.NewUUID()
	if err != nil {
		return "", 0, fmt.Errorf("generate new filename by uuid for save file: %w", err)
	}
	filename = id.String()

	dir := img.basePath + prefixPath + img.directoryToSave()
	err = os.MkdirAll(dir, 0750)
	if err != nil {
		return "", 0, fmt.Errorf("mkdir %s to save file: %w", dir, err)
	}

	filename = dir + filename + "." + extension
	file, err := os.Create(filename)
	if err != nil {
		return "", 0, fmt.Errorf("create %s to save file: %w", filename, err)
	}
	defer file.Close()

	written, err = io.Copy(file, image)
	return
}

func (img *imageRepoFS) SetBasePath(path string) {
	img.m.Lock()
	img.basePath = path
	img.m.Unlock()
}

func (img *imageRepoFS) SetDirToSave(fn func() string) {
	img.m.Lock()
	img.directoryToSave = fn
	img.m.Unlock()
}

func dirToSave() string {
	return time.Now().UTC().Format("2006/01/02/")
}
