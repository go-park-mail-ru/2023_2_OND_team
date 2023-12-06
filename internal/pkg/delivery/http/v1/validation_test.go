package v1

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFetchValidParams(t *testing.T) {
	rawUrl := "https://domain.test:8080/api/v1/pin"

	tests := []struct {
		name              string
		queryRow          string
		wantCount, lastID int
	}{
		{"both parameters were passed correctly", "?count=6&lastID=12", 6, 12},
		{"both parameters were passed correctly in a different order", "?lastID=22&count=3", 3, 22},
		{"repeating parameters", "?count=14&maxID=9&count=3&lastID=1", 14, 1},
		{"repeating parameters", "?count=14&lastID=1&count=3&mmmmmID=55&ID=65", 14, 1},
		{"empty parameters minID, maxID", "?count=7", 7, 0},
		{"the parameter maxID is registered but not specified", "?lastID=&count=17", 17, 0},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			URL, err := url.Parse(rawUrl + test.queryRow)
			if err != nil {
				t.Fatalf("error when parsing into the url.URL structure: %v", err)
			}
			actualCount, actualLastID, err := FetchValidParamForLoadFeed(URL)
			require.NoError(t, err, test.name)
			require.Equal(t, test.wantCount, actualCount, test.name)
			require.Equal(t, test.lastID, actualLastID, test.name)
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
		{"negative count", "?count=-5&lastID=12", ErrBadParams},
		{"negative ID", "?count=5&lastID=-6", ErrBadParams},
		{"requested count is more than a thousand", "?count=1001", ErrBadParams},
		{"count param empty", "?count=&lastID=23", ErrCountParameterMissing},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			URL, err := url.Parse(rawUrl + test.queryRow)
			if err != nil {
				t.Fatalf("error when parsing into the url.URL structure: %v", err)
			}
			actualCount, actualLastID, err := FetchValidParamForLoadFeed(URL)
			require.ErrorIs(t, err, test.wantErr)
			require.Equal(t, 0, actualCount)
			require.Equal(t, 0, actualLastID)
		})
	}
}
