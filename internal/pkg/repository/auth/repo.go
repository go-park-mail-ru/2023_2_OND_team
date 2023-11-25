package auth

import (
	"context"
	"fmt"

	authProto "github.com/go-park-mail-ru/2023_2_OND_team/api/auth"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/entity/session"
	entity "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/entity/user"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Repository interface {
	Register(ctx context.Context, user *entity.User) error
	Logout(ctx context.Context, sess *session.Session) error
	Login(ctx context.Context, username, password string) (*session.Session, error)
	GetUserID(ctx context.Context, sess *session.Session) (int, error)
}

type authRepo struct {
	client authProto.AuthClient
}

func NewAuthRepo(c authProto.AuthClient) *authRepo {
	return &authRepo{c}
}

func (r *authRepo) Register(ctx context.Context, user *entity.User) error {
	_, err := r.client.Register(ctx, &authProto.RegisterData{
		Cred: &authProto.Credentials{
			Username: user.Username,
			Password: user.Password,
		},
		Email: user.Email,
	})
	return err
}

func (r *authRepo) Logout(ctx context.Context, sess *session.Session) error {
	_, err := r.client.Logout(ctx, &authProto.Session{
		Key:    sess.Key,
		UserID: int64(sess.UserID),
		Expire: timestamppb.New(sess.Expire),
	})
	return err
}

func (r *authRepo) Login(ctx context.Context, username, password string) (*session.Session, error) {
	sess, err := r.client.Login(ctx, &authProto.Credentials{
		Username: username,
		Password: password,
	})
	if err != nil {
		return nil, fmt.Errorf("login: %w", err)
	}
	return &session.Session{
		Key:    sess.Key,
		UserID: int(sess.UserID),
		Expire: sess.Expire.AsTime(),
	}, nil
}

func (r *authRepo) GetUserID(ctx context.Context, sess *session.Session) (int, error) {
	userID, err := r.client.GetUserID(ctx, &authProto.Session{
		Key:    sess.Key,
		UserID: int64(sess.UserID),
		Expire: timestamppb.New(sess.Expire),
	})
	if err != nil {
		return 0, fmt.Errorf("get user id: %w", err)
	}
	return int(userID.Id), nil
}
