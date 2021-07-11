package goutils

import (
	"fmt"
	"testing"
)

func TestIsInteger(t *testing.T) {
	if ok := IsInteger("1"); !ok {
		t.Errorf("IsInteger 1 test failed.")
	}
	if ok := IsInteger("12"); !ok {
		t.Errorf("IsInteger 12 test failed.")
	}
	if ok := IsInteger("-123"); !ok {
		t.Errorf("IsInteger -123 test failed.")
	}
}

func TestIsRegularIPv4Address(t *testing.T) {
	if ok := IsRegularIPv4Address("192.168.1.147"); !ok {
		t.Errorf("IsRegularIPv4Address test failed.")
	}
	if ok := IsRegularIPv4Address("192.168.1.0/24"); !ok {
		t.Errorf("IsRegularIPv4Address test failed.")
	}
}

func TestPathNormalized(t *testing.T) {
	s := RemovePathSeparatorSuffix(PathNormalized("C:/tmp\\a////b/c\\\\"))

	fmt.Println(s)
}
