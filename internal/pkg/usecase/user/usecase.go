package user

import (
	"context"
	"errors"
	"io"

	entity "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/entity/user"
	repo "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/repository/user"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/usecase/image"
	"github.com/go-park-mail-ru/2023_2_OND_team/pkg/logger"
)

var ErrUserAuthentication = errors.New("user authentication")

const (
	lenSalt         = 16
	lenPasswordHash = 64
)

type Usecase interface {
	Register(ctx context.Context, user *entity.User) error
	Authentication(ctx context.Context, credentials userCredentials) (*entity.User, error)
	FindOutUsernameAndAvatar(ctx context.Context, userID int) (username string, avatar string, err error)
	UpdateUserAvatar(ctx context.Context, userID int, mimeTypeAvatar string, sizeAvatar int64, avatar io.Reader) error
	GetAllProfileInfo(ctx context.Context, userID int) (*entity.User, error)
	EditProfileInfo(ctx context.Context, userID int, updateData *profileUpdateData) error
}

type userCase struct {
	image.Usecase
	log  *logger.Logger
	repo repo.Repository
}

func New(log *logger.Logger, imgCase image.Usecase, repo repo.Repository) *userCase {
	return &userCase{
		Usecase: imgCase,
		log:     log,
		repo:    repo,
	}
}
