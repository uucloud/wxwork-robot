package utils

import "strings"

func TrimAt(name, text string) string {
	return strings.TrimSpace(strings.ReplaceAll(text, "@"+name, ""))
}
