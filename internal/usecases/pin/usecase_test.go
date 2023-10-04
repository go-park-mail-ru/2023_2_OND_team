package pin

import (
	"context"
	"math/rand"
	"strconv"
	"testing"

	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/repository/ramrepo"
	"github.com/go-park-mail-ru/2023_2_OND_team/pkg/logger"
	"github.com/stretchr/testify/require"
)

func TestSelectNewPins(t *testing.T) {
	log, _ := logger.New(logger.RFC3339FormatTime())
	defer log.Sync()

	db, _ := ramrepo.OpenDB(strconv.FormatInt(int64(rand.Int()), 10))
	defer db.Close()

	pinCase := New(log, ramrepo.NewRamPinRepo(db))

	testCases := []struct {
		name          string
		count, lastID int
		expNewLastID  int
	}{
		{
			name:         "provide correct count and lastID",
			count:        2,
			lastID:       1,
			expNewLastID: 3,
		},
		{
			name:         "provide incorrect count",
			count:        -2,
			lastID:       1,
			expNewLastID: 1,
		},
	}

	for _, tCase := range testCases {
		t.Run(tCase.name, func(t *testing.T) {
			_, actualLastID := pinCase.SelectNewPins(context.Background(), tCase.count, tCase.lastID)
			require.Equal(t, tCase.expNewLastID, actualLastID)
		})
	}
}
