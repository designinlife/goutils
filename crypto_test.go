package goutils

import (
	"testing"
)

func TestMD5(t *testing.T) {
	if ans := MD5("test"); ans != "098f6bcd4621d373cade4e832627b4f6" {
		t.Errorf("MD5 测试失败. (%s)", ans)
	}
}

func TestSHA1(t *testing.T) {
	if ans := SHA1("test"); ans != "a94a8fe5ccb19ba61c4c0873d391e987982fbbd3" {
		t.Errorf("SHA1 测试失败. (%s)", ans)
	}
}

func TestSHA2(t *testing.T) {
	if ans := SHA2("test"); ans != "90a3ed9e32b2aaf4c61c410eb925426119e1a9dc53d4286ade99a809" {
		t.Errorf("SHA224 测试失败. (%s)", ans)
	}
}

func TestSHA256(t *testing.T) {
	if ans := SHA256("test"); ans != "9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08" {
		t.Errorf("SHA256 测试失败. (%s)", ans)
	}
}

func TestSHA3(t *testing.T) {
	if ans := SHA3("test"); ans != "768412320f7b0aa5812fce428dc4706b3cae50e02a64caa16a782249bfe8efc4b7ef1ccb126255d196047dfedf17a0a9" {
		t.Errorf("SHA384 测试失败. (%s)", ans)
	}
}

func TestSHA512(t *testing.T) {
	if ans := SHA512("test"); ans != "ee26b0dd4af7e749aa1a8ee3c10ae9923f618980772e473f8819a5d4940e0db27ac185f8a0e1d5f84f88bc887fd67b143732c304cc5fa9ad8e6f57f50028a8ff" {
		t.Errorf("SHA512 测试失败. (%s)", ans)
	}
}
