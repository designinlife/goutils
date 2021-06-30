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
	"path/filepath"
)

// RemoveAny 移除任意文件或目录。（当传入的参数是系统保护路径时会报错！）
func RemoveAny(s string) error {
	// 保护系统根路径
	if InSlice([]string{"/", "/bin", "/boot", "/data", "/dev", "/etc", "/home", "/lib", "/lib64", "/media", "/mnt", "/opt", "/proc", "/root", "/run", "/sbin", "/srv", "/sys", "/tmp", "/usr", "/var"}, s) {
		return errors.New(fmt.Sprintf("不允许删除系统路径。（%s）", s))
	}

	err := os.RemoveAll(s)

	if err != nil {
		return err
	}

	return nil
}

// RemoveContents 移除目录下的所有文件及子目录。（不包含目录自身）
func RemoveContents(dir string) error {
	if !IsDir(dir) {
		return errors.New(fmt.Sprintf("参数必须是一个目录路径。(%s)", dir))
	}

	d, err := os.Open(dir)
	if err != nil {
		return err
	}

	defer d.Close()

	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}

	for _, name := range names {
		err = os.RemoveAll(filepath.Join(dir, name))
		if err != nil {
			return err
		}
	}

	return nil
}

// RemoveGlob 按模式匹配规则移除内容。
func RemoveGlob(path string) (err error) {
	contents, err := filepath.Glob(path)
	if err != nil {
		return
	}
	for _, item := range contents {
		err = os.RemoveAll(item)
		if err != nil {
			return
		}
	}
	return
}

// CheckSum 计算文件哈希校验码。
func CheckSum(filename string, algorithm string) (string, error) {
	f, err := os.Open(filename)

	if err != nil {
		return "", err
	}

	defer f.Close()

	var h hash.Hash

	switch algorithm {
	case "md5", "MD5":
		h = md5.New()
	case "sha1", "SHA1":
		h = sha1.New()
	case "sha224", "SHA224":
		h = sha256.New224()
	case "sha256", "SHA256":
		h = sha256.New()
	case "sha384", "SHA384":
		h = sha512.New384()
	case "sha512", "SHA512":
		h = sha512.New()
	default:
		return "", errors.New(fmt.Sprintf("不支持的 Hash (%s) 算法。", algorithm))
	}

	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	sum := h.Sum(nil)

	return fmt.Sprintf("%x", sum), nil
}

// IsFile 检查是否文件？
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

// IsDir 检查是否文件夹？
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
