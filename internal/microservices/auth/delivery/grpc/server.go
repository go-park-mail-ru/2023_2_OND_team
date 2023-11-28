package auth

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	authProto "github.com/go-park-mail-ru/2023_2_OND_team/internal/api/auth"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/entity/user"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/usecase/session"
	userUsecase "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/usecase/user"
	"github.com/go-park-mail-ru/2023_2_OND_team/pkg/logger"
)

type Usecase interface {
	Register(ctx context.Context, user *user.User) error
	Authentication(ctx context.Context, credentials userUsecase.UserCredentials) (*user.User, error)
}

type AuthServer struct {
	authProto.UnimplementedAuthServer

	log      *logger.Logger
	sm       session.SessionManager
	userCase Usecase
}

func New(log *logger.Logger, sm session.SessionManager, userCase Usecase) AuthServer {
	return AuthServer{
		UnimplementedAuthServer: authProto.UnimplementedAuthServer{},
		log:                     log,
		sm:                      sm,
		userCase:                userCase,
	}
}

func (as AuthServer) Register(ctx context.Context, rd *authProto.RegisterData) (*empty.Empty, error) {
	user := &user.User{
		Email:    rd.Email,
		Username: rd.Cred.Username,
		Password: rd.Cred.Password,
	}

	err := as.userCase.Register(ctx, user)
	if err != nil {
		as.log.Error(err.Error())
		return nil, status.Error(codes.Internal, "")
	}
	return &empty.Empty{}, nil
}

func (as AuthServer) Login(ctx context.Context, cred *authProto.Credentials) (*authProto.Session, error) {
	user, err := as.userCase.Authentication(ctx, userUsecase.UserCredentials{
		Username: cred.Username,
		Password: cred.Password,
	})
	if err != nil {
		as.log.Error(err.Error())
		return nil, status.Error(codes.Unauthenticated, "failed authentication")
	}

	session, err := as.sm.CreateNewSessionForUser(ctx, user.ID)
	if err != nil {
		as.log.Error(err.Error())
		return nil, status.Error(codes.Internal, "failed to create a session for the user")
	}

	return &authProto.Session{
		Key:    session.Key,
		UserID: int64(session.UserID),
		Expire: timestamppb.New(session.Expire),
	}, nil
}

func (as AuthServer) Logout(ctx context.Context, sess *authProto.Session) (*empty.Empty, error) {
	err := as.sm.DeleteUserSession(ctx, sess.Key)
	if err != nil {
		as.log.Error(err.Error())
		return nil, status.Error(codes.Internal, "delete user session")
	}
	return &empty.Empty{}, nil
}

func (as AuthServer) GetUserID(ctx context.Context, sess *authProto.Session) (*authProto.UserID, error) {
	userID, err := as.sm.GetUserIDBySessionKey(ctx, sess.Key)
	if err != nil {
		as.log.Error(err.Error())
		return nil, status.Error(codes.NotFound, "session not found")
	}
	return &authProto.UserID{Id: int64(userID)}, nil
}
