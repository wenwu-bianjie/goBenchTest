package util

import (
	"fmt"
	"testing"
)

func TestFormatStringWithBraceToSlice(t *testing.T) {
	var s = "{01000000, 01000001}"

	res := FormatStringWithBraceToSlice(s)

	fmt.Println(len(res))
}

func TestSubRegMatch(t *testing.T) {
	var s = " sub[ 1  , 2 ]"

	res := SubRegMatch(s)

	for _, m := range res {
		fmt.Println(string(m[0]))
		fmt.Println(string(m[1]))
		fmt.Println(string(m[2]))
	}
}
