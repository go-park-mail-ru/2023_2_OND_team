package validator

import "strings"

type (
	CheckRune func(r rune) bool
	CheckLen  func(length int) bool
)

func InRange(left, right int) CheckLen {
	return func(length int) bool {
		return left <= length && length <= right
	}
}

func Less(maxLength int) CheckLen {
	return func(length int) bool {
		return length <= maxLength
	}
}

func IsValidString(str string, validLen CheckLen, validRune ...CheckRune) bool {
	if !validLen(len(str)) {
		return false
	}

	for _, r := range str {
		for _, check := range validRune {
			if !check(r) {
				return false
			}
		}
	}
	return true
}

func IsValidListFromString(str, sep string, validLenEveryone CheckLen, maxLenList int,
	validRuneEveryone ...CheckRune) ([]string, bool) {

	list := strings.Split(str, sep)
	if len(list) > maxLenList {
		return nil, false
	}

	for _, part := range list {
		if !IsValidString(part, validLenEveryone, validRuneEveryone...) {
			return nil, false
		}
	}
	return list, true
}
