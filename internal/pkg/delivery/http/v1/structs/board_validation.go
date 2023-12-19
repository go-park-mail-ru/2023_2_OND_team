package structs

import (
	"fmt"
	"unicode"
)

func isValidTagTitle(title string) bool {
	if len(title) > 20 {
		return false
	}

	for _, sym := range title {
		if !(unicode.IsNumber(sym) || unicode.IsLetter(sym) || unicode.IsPunct(sym) || unicode.IsSpace(sym)) {
			return false
		}
	}
	return true
}

func checkIsValidTagTitles(titles []string) error {
	if len(titles) > 7 {
		return fmt.Errorf("too many titles")
	}

	invalidTitles := make([]string, 0)
	for _, title := range titles {
		if !isValidTagTitle(title) {
			invalidTitles = append(invalidTitles, title)
		}
	}
	if len(invalidTitles) > 0 {
		return fmt.Errorf("%v", invalidTitles)
	}
	return nil
}

func isValidBoardTitle(title string) bool {
	if len(title) == 0 || len(title) > 40 {
		return false
	}
	for _, sym := range title {
		if !(unicode.IsNumber(sym) || unicode.IsLetter(sym) || unicode.IsPunct(sym) || unicode.IsSpace(sym)) {
			return false
		}
	}
	return true
}
