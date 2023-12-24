package validation

import (
	"encoding/base64"
	"log"
	"strings"

	goaway "github.com/TwiN/go-away"
	"github.com/go-park-mail-ru/2023_2_OND_team/pkg/levenstein"
)

var additionalProfanity = []string{
	"0YXRg9C5",
	"0L/QuNC30LTQsA==",
	"0YfQu9C10L0=",
	"0LLQsNCz0LjQvdCw",
	"0LPQvtCy0L3Qvg==",
	"0L/QsNGA0LDRiNCw",
	"0YHRg9C60LA=",
	"0L/QuNGB0YzQutCw",
	"0YHQuNGB0YzQutC4",
	"0LLQu9Cw0LPQsNC70LjRidC1",
	"0L/QtdC90LjRgQ==",
	"0LHQu9GP0LTRjA==",
	"0YjQu9GO0YXQsA==",
	"0L/RgNC+0YHRgtC40YLRg9GC0LrQsA==",
	"0L3QuNCz0LXRgA==",
	"0L3QtdCz0YA=",
	"0YPQt9C60L7Qs9C70LDQt9GL0Lk=",
	"0YXQtdGA",
	"0LXQsdCw0YLRjA==",
	"0YPQtdCx0LDQvQ==",
	"0YPQtdCx0L7Qug==",
	"0YLRgNCw0YXQsNGC0Yw=",
	"0YLQstCw0YDRjA==",
	"0L/QuNC00YA=",
	"0YPRgNC+0LQ=",
	"0L/QuNC30LTQtdGG",
	"0YXRg9GP",
	"0LfQsNC70YPQv9Cw",
	"0L/QuNC00LDRgNCw0YE=",
	"0LvQvtGF",
	"0LPQsNC90LTQvtC9",
	"0LTRgNC+0YfQuNGC0Yw=",
	"0LDQvdCw0Ls=",
	"0LbQvtC/0LA=",
	"0LPQvdC40LTQsA==",
	"0YPRiNC70LXQv9C+0Lo=",
	"0YHRg9GH0LXQvdGL0Yg=",
	"0YHQv9C10YDQvNCw",
	"0LHQu9GP0YLRjA==",
	"0L/QvtGA0L3Qvg==",
	"0YHRgNCw0YLRjA==",
	"0YfQvNC+",
	"0LTQtdCx0LjQuw==",
	"0LrRgNC10YLQuNC9",
	"0LXQsdCw",
}

func GetLabels() []string {
	decodedLabels := make([]string, 0, len(additionalProfanity))
	for _, badEncoded := range additionalProfanity {
		decoded, err := base64.StdEncoding.DecodeString(badEncoded)
		if err != nil {
			log.Println(err)
		}
		decodedLabels = append(decodedLabels, string(decoded))
	}
	return decodedLabels
}

type ProfanityCensor interface {
	IsProfane(string) bool
	Sanitize(string) string
}

type defaultCensor struct {
	censor *goaway.ProfanityDetector
}

func isSimilarWithProfane(s string) bool {
	for _, badWord := range GetLabels() {
		if float64(levenstein.Levenshtein([]rune(s), []rune(badWord))) <= 0.3*float64(len(s)) {
			return true
		}
	}
	return false
}

func (c *defaultCensor) IsProfane(s string) bool {
	for _, word := range strings.Fields(s) {
		if c.censor.IsProfane(strings.ToLower(word)) || isSimilarWithProfane(strings.ToLower(word)) {
			return true
		}
	}
	return false
}

func (c *defaultCensor) Sanitize(s string) string {
	return c.censor.Censor(s)
}

func NewCensor(censor *goaway.ProfanityDetector) ProfanityCensor {
	return &defaultCensor{censor}
}
