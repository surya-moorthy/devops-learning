package looping

import "strings"

func Repeat(char string) string {
    var repeated strings.Builder
	for i:= 0; i<6 ; i++ {
		repeated.WriteString(char)
	}
	return repeated.String()
}