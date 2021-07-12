package goutils

import (
	"testing"
)

func TestZip2(t *testing.T) {
	err := Zip("D://tmp\\YX", "D:\\tmp\\test-zip.zip", false)

	if err != nil {
		t.Fatal(err)
	}
}

func TestUnzip(t *testing.T) {
	err := Unzip("D:\\tmp\\test2-zip", "D:\\tmp\\YX2")

	if err != nil {
		t.Fatal(err)
	}
}
