package service

import (
	"fmt"
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"
)

type TestCaseFetch struct {
	rawURL              string
	expCount, expLastID int
	expErr              error
}

func TestFetchValidParamForLoadTape(t *testing.T) {
	rawUrl := "https://domain.test:8080/api/v1/pin"
	testCases := make([]TestCaseFetch, 0)

	for count := 1; count != 5; count++ {
		for lastID := 1; lastID != 5; lastID++ {
			rawUrlWithParams := fmt.Sprintf("%s?count=%d&lastID=%d", rawUrl, count, lastID)
			testCases = append(testCases, TestCaseFetch{
				rawURL:    rawUrlWithParams,
				expCount:  count,
				expLastID: lastID,
				expErr:    nil,
			})
		}
	}
	testCases = append(testCases, []TestCaseFetch{
		{
			rawURL:    fmt.Sprintf("%s?count=%d&lastID=%d", rawUrl, 0, 1),
			expCount:  0,
			expLastID: 0,
			expErr:    ErrBadParams,
		},
		{
			rawURL:    fmt.Sprintf("%s?count=%d&lastID=%d", rawUrl, 12312, 1),
			expCount:  0,
			expLastID: 0,
			expErr:    ErrBadParams,
		},
		{
			rawURL:    fmt.Sprintf("%s?count=%d&lastID=%d", rawUrl, -2, 1),
			expCount:  0,
			expLastID: 0,
			expErr:    ErrBadParams,
		},
		{
			rawURL:    fmt.Sprintf("%s?count=%d&lastID=%d", rawUrl, 2, -1),
			expCount:  0,
			expLastID: 0,
			expErr:    ErrBadParams,
		},
		{
			rawURL:    fmt.Sprintf("%s?count=&lastID=%d", rawUrl, 1),
			expCount:  0,
			expLastID: 0,
			expErr:    ErrCountParameterMissing,
		},
		{
			rawURL:    fmt.Sprintf("%s?&lastID=%d", rawUrl, 4),
			expCount:  0,
			expLastID: 0,
			expErr:    ErrCountParameterMissing,
		},
	}...)

	for _, tCase := range testCases {
		t.Run("TestFetchValidParamForLoadTape", func(t *testing.T) {
			URL, _ := url.Parse(tCase.rawURL)
			actualCount, actualLastID, actualErr := FetchValidParamForLoadTape(URL)
			require.Equal(t, tCase.expCount, actualCount)
			require.Equal(t, tCase.expLastID, actualLastID)
			require.Equal(t, tCase.expErr, actualErr)
		})
	}
}
