package util

import (
	"regexp"
	"strings"
)

var reg = regexp.MustCompile("[^{^}^\\s^\\,]+")

var SubReg = regexp.MustCompile("sub\\[[\\s]*([0-9]+)[\\s]*\\,[\\s]*([0-9]+)[\\s]*\\]")

func SubStringFirstWord(s string) (res string) {
	res = string([]rune(s)[:1])
	return
}

func RemoveStringFirstWord(s string) (res string) {
	res = string([]rune(s)[1:])
	return
}

func FormatStringWithBraceToSlice(s string) (res []string) {
	var (
		m [][][]byte
	)
	m = reg.FindAllSubmatch([]byte(s), -1)

	for _, v := range m {
		res = append(res, string(v[0]))
	}
	return
}

func TrimStr(s string) string {
	return strings.Replace(s, " ", "", -1)
}

func SubRegMatch(s string) [][][]byte {
	res := SubReg.FindAllSubmatch([]byte(s), -1)
	return res
}
