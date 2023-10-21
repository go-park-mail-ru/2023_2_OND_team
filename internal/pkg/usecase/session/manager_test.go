package session

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"
	"testing"

	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/repository/ramrepo"
	"github.com/go-park-mail-ru/2023_2_OND_team/pkg/logger"
	"github.com/stretchr/testify/require"
)

func TestGetUserIDBySessionKey(t *testing.T) {
	log, err := logger.New(logger.RFC3339FormatTime())
	if err != nil {
		fmt.Println(err)
		return
	}
	defer log.Sync()

	db, err := ramrepo.OpenDB(strconv.FormatInt(int64(rand.Int()), 10))
	if err != nil {
		log.Error(err.Error())
		return
	}
	defer db.Close()

	sm := New(log, ramrepo.NewRamSessionRepo(db))

	testCases := []struct {
		name        string
		session_key string
		expUserId   int
		expErr      error
	}{
		{
			"providing valid session key",
			"461afabf38b3147c",
			1,
			nil,
		},
	}

	for _, tCase := range testCases {
		t.Run(tCase.name, func(t *testing.T) {
			id, err := sm.GetUserIDBySessionKey(context.Background(), tCase.session_key)
			require.Equal(t, tCase.expErr, err)
			require.Equal(t, tCase.expUserId, id)
		})
	}
}
