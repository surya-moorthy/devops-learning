package looping

import "strings"

func Repeat(char string,RepeadCount int) string {
    var repeated strings.Builder
	for i:= 0; i< RepeadCount ; i++ {
		repeated.WriteString(char)
	}
	return repeated.String()
}