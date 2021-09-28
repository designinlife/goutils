package goutils

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
	"github.com/pkg/errors"
	"hash"
	"io"
	"os"
	"path"
	"path/filepath"
)

// RemoveAllSafe 移除任意文件或目录。（当传入的参数是系统保护路径时会报错！）
func RemoveAllSafe(s string) error {
	// 保护系统根路径
	if InSlice([]string{"/", "/bin", "/boot", "/data", "/dev", "/etc", "/home", "/lib", "/lib64", "/media", "/mnt", "/opt", "/proc", "/root", "/run", "/sbin", "/srv", "/sys", "/tmp", "/usr", "/usr/bin", "/usr/sbin", "/usr/local/bin", "/usr/local/sbin", "/usr/local/etc", "/var"}, s) {
		return errors.New(fmt.Sprintf("It is strictly forbidden to delete the protected path. (%s)", s))
	}
	if InSlicePrefix([]string{"/bin", "/usr/bin", "/usr/sbin", "/etc", "/dev", "/lib", "/lib64", "/media", "/boot", "/proc", "/sbin", "/sys"}, s) {
		return errors.New(fmt.Sprintf("It is strictly forbidden to delete the protected path. (%s)", s))
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
		return errors.New(fmt.Sprintf("The parameter must be a directory. (%s)", dir))
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

// CopyFile 拷贝文件。
func CopyFile(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}

// VerifySum 校验文件哈希值。
func VerifySum(filename, checksum string, algorithm HashAlgorithm) bool {
	if !IsFile(filename) {
		return false
	}

	code, err := CheckSum(filename, algorithm, false)

	if err != nil {
		return false
	}

	if code == checksum {
		return true
	}

	return false
}

// CheckSum 计算文件哈希校验码。
func CheckSum(filename string, algorithm HashAlgorithm, capital bool) (string, error) {
	if !IsFile(filename) {
		return "", errors.New(fmt.Sprintf("File not found. (%s)", filename))
	}

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
		return "", errors.New(fmt.Sprintf("Unsupported hash algorithm. (%s)", algorithm))
	}

	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	sum := h.Sum(nil)

	if capital {
		return fmt.Sprintf("%X", sum), nil
	}

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

// SearchFile 在若干目录中搜索 name 文件。
func SearchFile(name string, dirs []string) (string, error) {
	var fn string

	for _, v := range dirs {
		fn = path.Join(v, name)

		if IsFile(fn) {
			return fn, nil
		}
	}

	return "", fmt.Errorf("no file was found. (%s)", name)
}
