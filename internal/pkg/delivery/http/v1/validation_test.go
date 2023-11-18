package v1

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFetchValidParams(t *testing.T) {
	rawUrl := "https://domain.test:8080/api/v1/pin"

	tests := []struct {
		name                            string
		queryRow                        string
		wantCount, wantMinID, wantMaxID int
	}{
		{"both parameters were passed correctly", "?count=6&minID=12&maxID=9", 6, 12, 9},
		{"both parameters were passed correctly in a different order", "?maxID=88&count=3&minID=22", 3, 22, 88},
		{"repeating parameters", "?count=14&maxID=9&count=3&maxID=55&minID=1", 14, 1, 9},
		{"repeating parameters", "?count=14&minID=1&count=3&mmmmmID=55&ID=65&maxID=1", 14, 1, 1},
		{"empty parameters minID, maxID", "?count=7", 7, 0, 0},
		{"the parameter maxID is registered but not specified", "?lastID=&count=17", 17, 0, 0},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			URL, err := url.Parse(rawUrl + test.queryRow)
			if err != nil {
				t.Fatalf("error when parsing into the url.URL structure: %v", err)
			}
			actualCount, actualMinID, actualMaxID, err := FetchValidParamForLoadTape(URL)
			require.NoError(t, err, test.name)
			require.Equal(t, test.wantCount, actualCount, test.name)
			require.Equal(t, test.wantMinID, actualMinID, test.name)
			require.Equal(t, test.wantMaxID, actualMaxID, test.name)
		})
	}
}

func TestErrorFetchValidParams(t *testing.T) {
	rawUrl := "https://domain.test:8080/api/v1/pin"

	tests := []struct {
		name     string
		queryRow string
		wantErr  error
	}{
		{"empty query row", "", ErrCountParameterMissing},
		{"count equal zero", "?count=0", ErrBadParams},
		{"negative count", "?count=-5&minID=12", ErrBadParams},
		{"negative ID", "?count=5&maxID=-6", ErrBadParams},
		{"requested count is more than a thousand", "?count=1001", ErrBadParams},
		{"count param empty", "?count=&minID=6&maxID=9", ErrCountParameterMissing},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			URL, err := url.Parse(rawUrl + test.queryRow)
			if err != nil {
				t.Fatalf("error when parsing into the url.URL structure: %v", err)
			}
			actualCount, actualLastID, _, err := FetchValidParamForLoadTape(URL)
			require.ErrorIs(t, err, test.wantErr)
			require.Equal(t, 0, actualCount)
			require.Equal(t, 0, actualLastID)
		})
	}
}
