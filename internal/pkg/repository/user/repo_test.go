package user

import (
	"context"
	"errors"
	"testing"

	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/entity/user"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/pashagolub/pgxmock/v2"
	"github.com/stretchr/testify/require"
)

func TestAddNewUser(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pool, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}

	repoUser := NewUserRepoPG(pool)

	pool.ExpectExec("INSERT INTO profile").
		WithArgs("my_username", "1234", "a@test.com").
		WillReturnResult(pgxmock.NewResult("INSERT", 1))

	err = repoUser.AddNewUser(ctx, &user.User{
		Username: "my_username",
		Password: "1234",
		Email:    "a@test.com",
	})
	require.NoError(t, err)

	wantErr := errors.New("insert profile fail")
	pool.ExpectExec("INSERT INTO profile").
		WithArgs("my_username", "1234", "a@test.com").
		WillReturnError(wantErr)

	err = repoUser.AddNewUser(ctx, &user.User{
		Username: "my_username",
		Password: "1234",
		Email:    "a@test.com",
	})
	require.ErrorIs(t, err, wantErr)
	require.EqualError(t, err, "add a new profile in storage: insert profile fail")
}

func TestGetUserByUsername(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pool, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}

	repoUser := NewUserRepoPG(pool)

	pool.ExpectQuery("SELECT id, password, email FROM profile").
		WithArgs("alex").
		WillReturnRows(
			pgxmock.NewRows([]string{"id", "password", "email"}).
				AddRow(3, "salt+hashPASS", "ex@mail.ru"),
		)
	wantUser := user.User{
		ID:       3,
		Password: "salt+hashPASS",
		Email:    "ex@mail.ru",
		Username: "alex",
	}
	actualUser, err := repoUser.GetUserByUsername(ctx, "alex")
	require.NoError(t, err)
	require.NotNil(t, actualUser)
	require.Equal(t, wantUser, *actualUser)

	wantErr := errors.New("fail select from profile")
	pool.ExpectQuery("SELECT id, password, email FROM profile").
		WithArgs("alex").
		WillReturnError(wantErr)
	actualUser, err = repoUser.GetUserByUsername(ctx, "alex")
	require.ErrorIs(t, err, wantErr)
	require.EqualError(t, err, "getting a user from storage: fail select from profile")
	require.Nil(t, actualUser)
}

func TestGetUsernameAndAvatarByID(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pool, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}

	repoUser := NewUserRepoPG(pool)

	pool.ExpectQuery("SELECT username, avatar FROM profile").
		WithArgs(14).
		WillReturnRows(
			pgxmock.NewRows([]string{"username", "avatar"}).
				AddRow("noname", "//pinspire.online/avatar.png"),
		)
	actualUsername, actualAvatar, err := repoUser.GetUsernameAndAvatarByID(ctx, 14)
	require.NoError(t, err)
	require.Equal(t, "noname", actualUsername)
	require.Equal(t, "//pinspire.online/avatar.png", actualAvatar)
}

func TestEditUserAvatar(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pool, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}

	repoUser := NewUserRepoPG(pool)

	pool.ExpectExec("UPDATE profile").
		WithArgs("new_avatar.svg", 16).
		WillReturnResult(pgxmock.NewResult("UPDATE", 1))
	err = repoUser.EditUserAvatar(ctx, 16, "new_avatar.svg")
	require.NoError(t, err)

	wantErr := errors.New("fail update")
	pool.ExpectExec("UPDATE profile").
		WithArgs("new_avatar.svg", 16).
		WillReturnError(wantErr)
	err = repoUser.EditUserAvatar(ctx, 16, "new_avatar.svg")
	require.ErrorIs(t, err, wantErr)
}

func TestGetAllUserData(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pool, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}

	repoUser := NewUserRepoPG(pool)

	pool.ExpectQuery("SELECT username, email, avatar, name, surname, about_me FROM profile").
		WithArgs(54).
		WillReturnRows(
			pgxmock.NewRows([]string{"username", "email", "avatar", "name", "surname", "about_me"}).
				AddRow("name", "12", "pic.webp", "", nil, "friendly"),
		)
	actualUser, err := repoUser.GetAllUserData(ctx, 54)
	require.NoError(t, err)
	require.NotNil(t, actualUser)
	require.Equal(t, user.User{
		ID:       54,
		Avatar:   "pic.webp",
		Surname:  pgtype.Text{String: "", Valid: false},
		Name:     pgtype.Text{String: "", Valid: true},
		Username: "name",
		AboutMe:  pgtype.Text{String: "friendly", Valid: true},
		Email:    "12",
	}, *actualUser)
}

func TestEditUserInfo(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pool, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}

	repoUser := NewUserRepoPG(pool)

	pool.ExpectExec("UPDATE profile").
		WithArgs("new", 16).
		WillReturnResult(pgxmock.NewResult("UPDATE", 1))
	err = repoUser.EditUserInfo(ctx, 16, S{"name": "new"})
	require.NoError(t, err)
}

func TestGetUserIdByUsername(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pool, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}

	repoUser := NewUserRepoPG(pool)

	pool.ExpectQuery("SELECT id FROM profile").
		WithArgs("uniqname").
		WillReturnRows(
			pgxmock.NewRows([]string{"id"}).
				AddRow(4),
		)
	id, err := repoUser.GetUserIdByUsername(ctx, "uniqname")
	require.NoError(t, err)
	require.Equal(t, 4, id)
}
