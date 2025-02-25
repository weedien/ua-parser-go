package uaparser

import (
	"regexp"
	"strings"
)

// ClientHints User Agent Client Hints data
type ClientHints struct {
	brands       []IBrand
	fullVerList  []IBrand
	mobile       bool
	model        string
	platform     string
	platformVer  string
	architecture string
	bitness      string
}

func NewClientHints(headers map[string]string) ClientHints {

	for k, v := range headers {
		headers[strings.ToLower(k)] = v
		delete(headers, k)
	}

	stripQuotes := func(str string) string {
		return strings.ReplaceAll(str, "\"", "")
	}

	itemListToArray := func(header string) []IBrand {
		if header == "" {
			return nil
		}
		var arr []IBrand

		tokens := strings.Split(regexp.MustCompile(`\\?"`).ReplaceAllString(header, ""), ",")
		for _, token := range tokens {
			if strings.Contains(token, ";") {
				parts := strings.Split(strings.Trim(token, " "), ";v=")
				arr = append(arr, IBrand{
					Name:    parts[0],
					Version: parts[1],
				})
			} else {
				arr = append(arr, IBrand{
					Name: token,
				})
			}
		}
		return arr
	}

	mobile, _ := regexp.MatchString(`\?1`, headers[CHHeaderMobile])

	return ClientHints{
		brands:       itemListToArray(headers[CHHeader]),
		fullVerList:  itemListToArray(headers[CHHeaderFullVerList]),
		mobile:       mobile,
		model:        stripQuotes(headers[CHHeaderModel]),
		platform:     stripQuotes(headers[CHHeaderPlatform]),
		platformVer:  stripQuotes(headers[CHHeaderPlatformVer]),
		architecture: stripQuotes(headers[CHHeaderArch]),
		bitness:      stripQuotes(headers[CHHeaderBitness]),
	}
}
