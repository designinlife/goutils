package goutils

import (
	"testing"
)

func TestIsDir(t *testing.T) {
	if ok := IsDir("/tmp"); !ok {
		t.Fatal("test IsDir failed.")
	}
}
