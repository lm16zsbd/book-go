package utl

import "strings"

func FirstLang(val string) string {
	if val == "" {
		return ""
	}
	idx := strings.Index(val, ",,")
	if idx == -1 {
		return val
	}
	return val[:idx]
}
