package goutils

import (
	"regexp"
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

// InSlicePrefix 检查字符串前缀是否在列表中？
func InSlicePrefix(s []string, prefix string) bool {
	for _, v := range s {
		if strings.HasPrefix(v, prefix) {
			return true
		}
	}

	return false
}

// IsRegularIPv4AndPort 检查字符串是否匹配 IP:Port 模式？
func IsRegularIPv4AndPort(s string) bool {
	if ok, _ := regexp.MatchString("^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5]):[0-9]+$", s); ok {
		return true
	}

	return false
}
