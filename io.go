package goutils

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"errors"
	"fmt"
	"github.com/dustin/go-humanize"
	"hash"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// RemoveAnySafe 移除任意文件或目录。（当传入的参数是系统保护路径时会报错！）
func RemoveAnySafe(s string) error {
	// 保护系统根路径
	if InSlice([]string{"/", "/bin", "/boot", "/data", "/dev", "/etc", "/home", "/lib", "/lib64", "/media", "/mnt", "/opt", "/proc", "/root", "/run", "/sbin", "/srv", "/sys", "/tmp", "/usr", "/usr/bin", "/usr/sbin", "/usr/local/bin", "/usr/local/sbin", "/usr/local/etc", "/var"}, s) {
		return errors.New(fmt.Sprintf("不允许删除系统路径。（%s）", s))
	}
	if InSlicePrefix([]string{"/bin", "/usr/bin", "/usr/sbin", "/etc", "/dev", "/lib", "/lib64", "/media", "/boot", "/proc", "/sbin", "/sys"}, s) {
		return errors.New(fmt.Sprintf("不允许删除此前缀的系统路径。（%s）", s))
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

// WriteCounter counts the number of bytes written to it. It implements to the io.Writer interface
// and we can pass this into io.TeeReader() which will report progress on each write cycle.
type WriteCounter struct {
	Total uint64
}

func (wc *WriteCounter) Write(p []byte) (int, error) {
	n := len(p)
	wc.Total += uint64(n)
	wc.PrintProgress()
	return n, nil
}

func (wc WriteCounter) PrintProgress() {
	// Clear the line by using a character return to go back to the start and remove
	// the remaining characters by filling it with spaces
	fmt.Printf("\r%s", strings.Repeat(" ", 35))

	// Return again and print current status of download
	// We use the humanize package to print the bytes in a meaningful way (e.g. 10 MB)
	fmt.Printf("\rDownloading... %s complete", humanize.Bytes(wc.Total))
}
