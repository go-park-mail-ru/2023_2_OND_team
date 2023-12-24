package validation

import "github.com/microcosm-cc/bluemonday"

type SanitizerXSS interface {
	Sanitize(string) string
}

type bluemondaySanitizer struct {
	sanitizer *bluemonday.Policy
}

func (san *bluemondaySanitizer) Sanitize(s string) string {
	return san.sanitizer.Sanitize(s)
}

func NewSanitizerXSS(sanitizer *bluemonday.Policy) SanitizerXSS {
	return &bluemondaySanitizer{sanitizer}
}
