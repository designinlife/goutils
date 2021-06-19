package goutils

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"errors"
	"fmt"
	"hash"
	"io"
	"os"
	"strings"
)

// Substitute 替换字符串。
func Substitute(source string, replacements map[string]string) string {
	for k, v := range replacements {
		source = strings.ReplaceAll(source, "{"+k+"}", v)
	}

	return source
}

func FileChecksum(filename string, algorithm string) (string, error) {
	f, err := os.Open(filename)

	if err != nil {
		return "", err
	}

	defer f.Close()

	var h hash.Hash

	switch algorithm {
	case "md5":
		h = md5.New()
		break
	case "sha", "sha1":
		h = sha1.New()
		break
	case "sha256":
		h = sha256.New()
		break
	case "sha512":
		h = sha512.New()
		break
	default:
		return "", errors.New(fmt.Sprintf("不支持的 Hash (%s) 算法。", algorithm))
	}

	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	sum := h.Sum(nil)

	return fmt.Sprintf("%x", sum), nil
}

func IsFile(filename string) bool {
	info, err := os.Stat(filename)

	if err != nil {
		return false
	}

	if os.IsNotExist(err) {
		return false
	}

	return !info.IsDir()
}

func IsDir(dirname string) bool {
	info, err := os.Stat(dirname)

	if err != nil {
		return false
	}

	if os.IsNotExist(err) {
		return false
	}

	return info.IsDir()
}
