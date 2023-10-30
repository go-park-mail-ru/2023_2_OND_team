package board

import "unicode"

func isValidTagTitle(title string) bool {
	for _, sym := range title {
		if !(unicode.IsNumber(sym) || unicode.IsLetter(sym)) {
			return false
		}
	}
	return true
}

func isValidTagTitles(titles []string) bool {
	if len(titles) > 7 {
		return false
	}
	for _, title := range titles {
		if !isValidTagTitle(title) {
			return false
		}
	}
	return true
}

func isValidBoardTitle(title string) bool {
	if len(title) < 4 || len(title) > 50 {
		return false
	}
	for _, sym := range title {
		if !(unicode.IsNumber(sym) || unicode.IsLetter(sym)) {
			return false
		}
	}
	return true
}
