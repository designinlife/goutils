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

func TestTar(t *testing.T) {
	err := Tar("D:\\tmp\\YX", "D:\\tmp\\")

	if err != nil {
		t.Fatal(err)
	}
}

func TestUntar(t *testing.T) {
	err := Untar("D:\\tmp\\YX.tar", "D:\\tmp\\YX-tar")

	if err != nil {
		t.Fatal(err)
	}
}

func TestGzip(t *testing.T) {
	err := Gzip("D:\\tmp\\YX", "D:\\tmp\\")

	if err != nil {
		t.Fatal(err)
	}
}

func TestTgz(t *testing.T) {
	err := TarGz("D:\\tmp\\YX", "D:\\tmp\\")

	if err != nil {
		t.Fatal(err)
	}
}

func TestUntarGz(t *testing.T) {
	err := UntarGz("D:\\tmp\\YX.tar.gz", "D:\\tmp\\YX-tar-gz")

	if err != nil {
		t.Fatal(err)
	}
}
