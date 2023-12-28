package validation

import (
	"encoding/base64"
	"log"
	"strings"

	goaway "github.com/TwiN/go-away"
)

var additionalProfanity = []string{
	"0YXRgw==",
	"0L/QuNC30LTQsA==",
	"0YfQu9C10L0=",
	"0LLQsNCz0LjQvQ==",
	"0LPQvtCy0L3Qvg==",
	"0L/QsNGA0LDRiNCw",
	"0YHRg9C60LA=",
	"0L/QuNGB0YzQug==",
	"0YHQuNGB0YzQug==",
	"0LLQu9Cw0LPQsNC70LjRiQ==",
	"0L/QtdC90LjRgQ==",
	"0LHQu9GP0LTRjA==",
	"0YjQu9GO0YU=",
	"0L/RgNC+0YHRgtC40YLRg9GC0Lo=",
	"0L3QuNCz0LXRgA==",
	"0L3QtdCz0YA=",
	"0YPQt9C60L7Qs9C70LDQt9GL",
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
	"0LfQsNC70YPQvw==",
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
	"cGl6ZGE=",
	"0LXQsdCw0Ls=",
	"0LPQsNCy0L3Qvg==",
	"0LPQvtC90LTQvtC9",
	"0YXRg9C1",
	"0LXQsdCw0L0=",
	"0LXQsdC70LDQvQ==",
	"0LXQsdGD0YfQuA==",
	"0LXQsdC70LjQstGL",
	"0L/QuNC00YDQuNC7",
	"0L/QvtGA0L3Rg9GF0LA=",
	"0LXQsdC70Y8=",
	"0YPQtdCx0LjRidC90Ys=",
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

func (c *defaultCensor) IsProfane(s string) bool {
	for _, word := range strings.Fields(s) {
		if c.censor.IsProfane(strings.ToLower(word)) {
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
