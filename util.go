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

// RemovePathSeparatorPrefix 移除开头的 /,\ 字符。
func RemovePathSeparatorPrefix(s string) string {
	return strings.TrimPrefix(strings.TrimPrefix(s, "/"), "\\")
}

// RemovePathSeparatorSuffix 移除尾部的 /,\ 字符。
func RemovePathSeparatorSuffix(s string) string {
	return strings.TrimSuffix(strings.TrimSuffix(s, "/"), "\\")
}

// IsInteger 检查字符串是否匹配整数模式？
func IsInteger(s string) bool {
	if ok, _ := regexp.MatchString("^[-]?[1-9]+[0-9]*$", s); ok {
		return true
	}

	return false
}

// IsRegularIPv4Address 检查字符串是否匹配 IPv4 模式？
func IsRegularIPv4Address(s string) bool {
	if ok, _ := regexp.MatchString("^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])(/[0-9]{1,2})?$", s); ok {
		return true
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

// IsRegularEmailAddress 检查字符串是否匹配电子邮件地址模式？
func IsRegularEmailAddress(s string) bool {
	if ok, _ := regexp.MatchString("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$", s); ok {
		return true
	}

	return false
}

// IsRegularPhoneNumber 检查字符串是否匹配电话号码模式？
func IsRegularPhoneNumber(s string) bool {
	if ok, _ := regexp.MatchString("^(?:(?:\\+?1\\s*(?:[.-]\\s*)?)?(?:\\(\\s*([2-9]1[02-9]|[2-9][02-8]1|[2-9][02-8][02-9])\\s*\\)|([2-9]1[02-9]|[2-9][02-8]1|[2-9][02-8][02-9]))\\s*(?:[.-]\\s*)?)?([2-9]1[02-9]|[2-9][02-9]1|[2-9][02-9]{2})\\s*(?:[.-]\\s*)?([0-9]{4})(?:\\s*(?:#|x\\.?|ext\\.?|extension)\\s*(\\d+))?$", s); ok {
		return true
	}

	return false
}
