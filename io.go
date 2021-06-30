package goutils

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
	"github.com/dustin/go-humanize"
	"github.com/pkg/errors"
	"hash"
	"io"
	"os"
	"path/filepath"
)

// RemoveAnySafe 移除任意文件或目录。（当传入的参数是系统保护路径时会报错！）
func RemoveAnySafe(s string) error {
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

// WriteCounter counts the number of bytes written to it. It implements to the io.Writer interface
// and we can pass this into io.TeeReader() which will report progress on each write cycle.
type WriteCounter struct {
	LoadedBytes        uint64
	TotalBytes         uint64
	ProgressBar        bool
	OnlyShowPercentage bool
}

func (w *WriteCounter) Write(p []byte) (int, error) {
	n := len(p)
	w.LoadedBytes += uint64(n)

	if w.ProgressBar {
		w.PrintProgress()
	}

	return n, nil
}

func (w WriteCounter) PrintProgress() {
	// Clear the line by using a character return to go back to the start and remove
	// the remaining characters by filling it with spaces
	// fmt.Printf("\r%s", strings.Repeat(" ", 35))

	// Return again and print current status of download
	// We use the humanize package to print the bytes in a meaningful way (e.g. 10 MB)
	if w.TotalBytes > 0 {
		if w.OnlyShowPercentage {
			fmt.Printf("\r%.2f%%", float64(w.LoadedBytes)*100.00/float64(w.TotalBytes))
		} else {
			fmt.Printf("\rDownloading... %s of %s complete", humanize.Bytes(w.LoadedBytes), humanize.Bytes(w.TotalBytes))
		}
	} else {
		fmt.Printf("\rDownloading... %s complete", humanize.Bytes(w.LoadedBytes))
	}
}
