package goutils

import (
	"testing"
	"time"
)

func TestSubProcess_ShellExec(t *testing.T) {
	subp := NewSubProcess()
	exitCode, err := subp.ShellExec("python d:/t.py", "python d:/t.py")
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("Exit Code: %d", exitCode)
}

func TestSubProcess_Run(t *testing.T) {
	subp := NewSubProcess()
	exitCode, err := subp.Run("python", "d:/t.py")
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("Exit Code: %d", exitCode)
}

func TestSubProcess_RunWithOption(t *testing.T) {
	subp := NewSubProcessWithOptions(ProcessOptionWithTimeout(15 * time.Second))
	exitCode, err := subp.Run("python", "d:/t.py")
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("Exit Code: %d", exitCode)
}
