package user

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"
	"testing"

	entity "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/entity/user"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/repository/ramrepo"
	"github.com/go-park-mail-ru/2023_2_OND_team/pkg/logger"
	"github.com/stretchr/testify/require"
)

func TestRegister(t *testing.T) {
	log, err := logger.New(logger.RFC3339FormatTime())
	if err != nil {
		fmt.Println(err)
		return
	}
	defer log.Sync()

	db, _ := ramrepo.OpenDB(strconv.FormatInt(int64(rand.Int()), 10))
	defer db.Close()

	userCase := New(log, ramrepo.NewRamUserRepo(db))

	testCases := []struct {
		name   string
		user   *entity.User
		expErr error
	}{
		{
			"providing valid user",
			&entity.User{
				Username: "valid_user",
				Password: "helloworld",
				Email:    "gggg@yandex.ru",
			},
			nil,
		},
	}

	for _, tCase := range testCases {
		t.Run(tCase.name, func(t *testing.T) {
			err := userCase.Register(context.Background(), tCase.user)
			require.Equal(t, tCase.expErr, err)
		})
	}
}
