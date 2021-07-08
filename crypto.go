package goutils

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"hash"
	"math/big"
	"net"
	"os"
	"path/filepath"
	"time"
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

// GenerateSelfSignedCertKey 生成自签名证书。
func GenerateSelfSignedCertKey(outdir string, keySize int, expire time.Duration, host string, alternateIPs []net.IP, alternateDNS []string) error {
	// 1.生成密钥对
	priv, err := rsa.GenerateKey(rand.Reader, keySize)

	if err != nil {
		return err
	}

	// 2.创建证书模板
	template := x509.Certificate{
		SerialNumber: big.NewInt(1), // 该号码表示CA颁发的唯一序列号，在此使用一个数来代表
		Issuer:       pkix.Name{},
		Subject: pkix.Name{
			CommonName: fmt.Sprintf("%s@%d", host, time.Now().Unix()),
		},
		NotBefore:   time.Now(),
		NotAfter:    time.Now().Add(expire),
		KeyUsage:    x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign, // 表示该证书是用来做服务端认证的
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
	}

	if ip := net.ParseIP(host); ip != nil {
		template.IPAddresses = append(template.IPAddresses, ip)
	} else {
		template.DNSNames = append(template.DNSNames, host)
	}

	template.IPAddresses = append(template.IPAddresses, alternateIPs...)
	template.DNSNames = append(template.DNSNames, alternateDNS...)

	// 3.创建证书,这里第二个参数和第三个参数相同则表示该证书为自签证书，返回值为DER编码的证书
	certificate, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	if err != nil {
		return err
	}

	// 4.将得到的证书放入pem.Block结构体中
	block := pem.Block{
		Type:    "CERTIFICATE",
		Headers: nil,
		Bytes:   certificate,
	}

	// 5.通过pem编码并写入磁盘文件
	file, err := os.Create(filepath.Join(outdir, "ca.crt"))
	if err != nil {
		return err
	}
	defer file.Close()
	pem.Encode(file, &block)

	// 6.将私钥中的密钥对放入pem.Block结构体中
	block = pem.Block{
		Type:    "RSA PRIVATE KEY",
		Headers: nil,
		Bytes:   x509.MarshalPKCS1PrivateKey(priv),
	}

	// 7.通过pem编码并写入磁盘文件
	file, err = os.Create(filepath.Join(outdir, "ca.key"))
	if err != nil {
		return err
	}
	defer file.Close()
	pem.Encode(file, &block)

	return nil
}
