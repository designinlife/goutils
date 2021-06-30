package goutils

import (
	"strings"
)

// Substitute 替换字符串。
func Substitute(source string, replacements map[string]string) string {
	for k, v := range replacements {
		source = strings.ReplaceAll(source, "{"+k+"}", v)
	}

	return source
}

// InSlice 检查字符串是否在列表中？
func InSlice(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}
