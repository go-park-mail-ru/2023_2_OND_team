package service

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFetchValidParams(t *testing.T) {
	rawUrl := "https://domain.test:8080/api/v1/pin"
	// testCases := make([]TestCaseFetch, 0)

	tests := []struct {
		name                  string
		queryRow              string
		wantCount, wantLastID int
	}{
		{"both parameters were passed correctly", "?count=6&lastID=12", 6, 12},
		{"both parameters were passed correctly in a different order", "?lastID=1&count=3", 3, 1},
		{"repeating parameters", "?count=14&lastID=1&count=3&lastID=55&lastID=65", 14, 1},
		{"repeating parameters", "?count=14&lastID=1&count=3&lastID=55&lastID=65", 14, 1},
		{"empty parameter lastID", "?count=7", 7, 0},
		{"the parameter lastID is registered but not specified", "?lastID=&count=17", 17, 0},
	}

	// testCases = append(testCases, []TestCaseFetch{
	// 	{
	// 		rawURL:    fmt.Sprintf("%s?count=%d&lastID=%d", rawUrl, 0, 1),
	// 		expCount:  0,
	// 		expLastID: 0,
	// 		expErr:    ErrBadParams,
	// 	},
	// 	{
	// 		rawURL:    fmt.Sprintf("%s?count=%d&lastID=%d", rawUrl, 12312, 1),
	// 		expCount:  0,
	// 		expLastID: 0,
	// 		expErr:    ErrBadParams,
	// 	},
	// 	{
	// 		rawURL:    fmt.Sprintf("%s?count=%d&lastID=%d", rawUrl, -2, 1),
	// 		expCount:  0,
	// 		expLastID: 0,
	// 		expErr:    ErrBadParams,
	// 	},
	// 	{
	// 		rawURL:    fmt.Sprintf("%s?count=%d&lastID=%d", rawUrl, 2, -1),
	// 		expCount:  0,
	// 		expLastID: 0,
	// 		expErr:    ErrBadParams,
	// 	},
	// 	{
	// 		rawURL:    fmt.Sprintf("%s?count=&lastID=%d", rawUrl, 1),
	// 		expCount:  0,
	// 		expLastID: 0,
	// 		expErr:    ErrCountParameterMissing,
	// 	},
	// 	{
	// 		rawURL:    fmt.Sprintf("%s?&lastID=%d", rawUrl, 4),
	// 		expCount:  0,
	// 		expLastID: 0,
	// 		expErr:    ErrCountParameterMissing,
	// 	},
	// }...)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			URL, err := url.Parse(rawUrl + test.queryRow)
			if err != nil {
				t.Fatalf("error when parsing into the url.URL structure: %v", err)
			}
			actualCount, actualLastID, err := FetchValidParamForLoadTape(URL)
			require.NoError(t, err)
			require.Equal(t, test.wantCount, actualCount)
			require.Equal(t, test.wantLastID, actualLastID)
		})
	}
}
