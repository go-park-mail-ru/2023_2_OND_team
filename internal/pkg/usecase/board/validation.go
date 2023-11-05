package board

import (
	"fmt"
	"unicode"
)

func (bCase *BoardUsecase) isValidTagTitle(title string) bool {
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

func (bCase *BoardUsecase) checkIsValidTagTitles(titles []string) error {
	if len(titles) > 7 {
		return fmt.Errorf("too many titles")
	}

	invalidTitles := make([]string, 0)
	for _, title := range titles {
		if !bCase.isValidTagTitle(title) {
			invalidTitles = append(invalidTitles, title)
		}
	}
	if len(invalidTitles) > 0 {
		return fmt.Errorf("%v", invalidTitles)
	}
	return nil
}

func (bCase *BoardUsecase) isValidBoardTitle(title string) bool {
	if len(title) == 0 || len(title) > 40 {
		return false
	}
	for _, sym := range title {
		if !(unicode.IsNumber(sym) || unicode.IsLetter(sym) || unicode.IsPunct(sym) || unicode.IsSpace(sym)) {
			return false
		}
	}
	bCase.sanitizer.Sanitize(title)
	return true
}

func (bCase *BoardUsecase) isValidUsername(username string) bool {
	if len(username) < 4 || len(username) > 50 {
		return false
	}
	for _, r := range username {
		if !(unicode.IsNumber(r) || unicode.IsLetter(r)) {
			return false
		}
	}
	bCase.sanitizer.Sanitize(username)

	return true
}
