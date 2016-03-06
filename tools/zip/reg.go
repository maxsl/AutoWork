package zip

import (
	"regexp"
	"strings"
)

func match(str string, regexpstr []*regexp.Regexp) bool {
	for _, v := range regexpstr {
		if !v.Match([]byte(str)) {
			continue
		}
		return true
	}
	return false
}

func getreg(regexpstr []string) []*regexp.Regexp {
	list := make([]*regexp.Regexp, 0, len(regexpstr))
	for _, v := range regexpstr {
		if !strings.HasPrefix(v, "^") {
			v = "^" + v
		}
		reg, err := regexp.Compile(v)
		if err != nil {
			continue
		}
		list = append(list, reg)
	}
	return list
}
