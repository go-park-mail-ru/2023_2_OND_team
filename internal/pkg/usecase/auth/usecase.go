package auth

import (
	"context"
	"fmt"

	authProto "github.com/go-park-mail-ru/2023_2_OND_team/internal/api/auth"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/entity/session"
	entity "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/entity/user"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Usecase interface {
	Register(ctx context.Context, user *entity.User) error
	Login(ctx context.Context, username, password string) (*session.Session, error)
	GetUserIDBySession(ctx context.Context, sess *session.Session) (int, error)
	Logout(ctx context.Context, sess *session.Session) error
}

type authCase struct {
	client authProto.AuthClient
}

func New(client authProto.AuthClient) *authCase {
	return &authCase{client}
}

func (ac *authCase) Register(ctx context.Context, user *entity.User) error {
	_, err := ac.client.Register(ctx, &authProto.RegisterData{
		Cred: &authProto.Credentials{
			Username: user.Username,
			Password: user.Password,
		},
		Email: user.Email,
	})
	if err != nil {
		return fmt.Errorf("register: %w", err)
	}
	return nil
}

func (ac *authCase) Logout(ctx context.Context, sess *session.Session) error {
	_, err := ac.client.Logout(ctx, &authProto.Session{
		Key:    sess.Key,
		UserID: int64(sess.UserID),
		Expire: timestamppb.New(sess.Expire),
	})
	if err != nil {
		return fmt.Errorf("logout: %w", err)
	}
	return nil
}

func (ac *authCase) Login(ctx context.Context, username, password string) (*session.Session, error) {
	sess, err := ac.client.Login(ctx, &authProto.Credentials{
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

func (ac *authCase) GetUserIDBySession(ctx context.Context, sess *session.Session) (int, error) {
	userID, err := ac.client.GetUserID(ctx, &authProto.Session{
		Key:    sess.Key,
		UserID: int64(sess.UserID),
		Expire: timestamppb.New(sess.Expire),
	})
	if err != nil {
		return 0, fmt.Errorf("get user id by session: %w", err)
	}
	return int(userID.Id), nil
}
