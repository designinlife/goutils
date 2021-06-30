package goutils

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
	"hash"
)

type HashAlgorithm string

const (
	Md5    HashAlgorithm = "md5"
	Sha1   HashAlgorithm = "sha1"
	Sha224 HashAlgorithm = "sha224"
	Sha256 HashAlgorithm = "sha256"
	Sha384 HashAlgorithm = "sha384"
	Sha512 HashAlgorithm = "sha512"
)

func MD5(s string) string {
	return Hash(s, Md5, false)
}

func SHA1(s string) string {
	return Hash(s, Sha1, false)
}

func SHA2(s string) string {
	return Hash(s, Sha224, false)
}

func SHA256(s string) string {
	return Hash(s, Sha256, false)
}

func SHA3(s string) string {
	return Hash(s, Sha384, false)
}

func SHA512(s string) string {
	return Hash(s, Sha512, false)
}

func Hash(s string, algorithm HashAlgorithm, capital bool) string {
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
		return ""
	}

	h.Write([]byte(s))

	if capital {
		return fmt.Sprintf("%X", h.Sum(nil))
	}

	return fmt.Sprintf("%x", h.Sum(nil))
}

func HMD5(s, key string) string {
	return HMAC(s, key, Md5, false)
}

func HSHA1(s, key string) string {
	return HMAC(s, key, Sha1, false)
}

func HSHA2(s, key string) string {
	return HMAC(s, key, Sha224, false)
}

func HSHA256(s, key string) string {
	return HMAC(s, key, Sha256, false)
}

func HSHA3(s, key string) string {
	return HMAC(s, key, Sha384, false)
}

func HSHA512(s, key string) string {
	return HMAC(s, key, Sha512, false)
}

func HMAC(str, key string, algorithm HashAlgorithm, capital bool) string {
	var h hash.Hash

	switch algorithm {
	case "md5", "MD5":
		h = hmac.New(md5.New, []byte(key))
	case "sha1", "SHA1":
		h = hmac.New(sha1.New, []byte(key))
	case "sha224", "SHA224":
		h = hmac.New(sha256.New224, []byte(key))
	case "sha256", "SHA256":
		h = hmac.New(sha256.New, []byte(key))
	case "sha384", "SHA384":
		h = hmac.New(sha512.New384, []byte(key))
	case "sha512", "SHA512":
		h = hmac.New(sha512.New, []byte(key))
	default:
		return ""
	}

	h.Write([]byte(str))

	if capital {
		return fmt.Sprintf("%X", h.Sum(nil))
	}

	return fmt.Sprintf("%x", h.Sum(nil))
}
