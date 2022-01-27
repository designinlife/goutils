package goutils

import (
	"bufio"
	"context"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"syscall"
	"time"
)

type ProcessOutLineHandle func(line string)

type ProcessOption func(*SubProcess)

func ProcessOptionWithTimeout(timeout time.Duration) ProcessOption {
	return func(c *SubProcess) {
		c.Option.Timeout = timeout
	}
}

func ProcessOptionWithWriter(w io.Writer) ProcessOption {
	return func(c *SubProcess) {
		c.Option.Writer = w
	}
}

type SubProcessOption struct {
	Debug      bool
	Quiet      bool
	Timeout    time.Duration
	HandleFunc ProcessOutLineHandle
	Shell      bool
	Env        []string
	Writer     io.Writer
}

type SubProcess struct {
	Option        *SubProcessOption
	isWindows     bool
	shellBin      string
	shellBinParam string
}

func DefaultSubProcessOption() *SubProcessOption {
	opt := &SubProcessOption{
		Debug:   false,
		Quiet:   false,
		Timeout: 120 * time.Second,
		Shell:   true,
	}

	return opt
}

func NewSubProcessWithOptions(opt ...ProcessOption) *SubProcess {
	s := &SubProcess{
		Option:    DefaultSubProcessOption(),
		isWindows: runtime.GOOS == "windows",
	}

	for _, op := range opt {
		op(s)
	}

	if s.isWindows {
		s.shellBin = "cmd.exe"
		s.shellBinParam = "/C"
	} else {
		s.shellBin = "/bin/sh"
		s.shellBinParam = "-c"
	}

	return s
}

func NewSubProcess() *SubProcess {
	s := &SubProcess{
		Option:    DefaultSubProcessOption(),
		isWindows: runtime.GOOS == "windows",
	}

	if s.isWindows {
		s.shellBin = "cmd.exe"
		s.shellBinParam = "/C"
	} else {
		s.shellBin = "/bin/sh"
		s.shellBinParam = "-c"
	}

	return s
}

func (s *SubProcess) ShellExec(arg ...string) (int, error) {
	if s.shellBin == "" || s.shellBinParam == "" {
		return -1, errors.New("Please use NewSubProcess() or NewSubProcessWithOptions().")
	}

	arg2 := []string{s.shellBinParam, strings.Join(arg, " && ")}

	return s.Run(s.shellBin, arg2...)
}

// Run 执行系统命令。
func (s *SubProcess) Run(command string, args ...string) (int, error) {
	if s.Option == nil {
		s.Option = DefaultSubProcessOption()
	}

	ctx, cancel := context.WithTimeout(context.Background(), s.Option.Timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, command, args...)

	if s.Option.Env != nil {
		cmd.Env = append(os.Environ(), s.Option.Env...)
	}

	stderr, _ := cmd.StderrPipe()
	stdout, _ := cmd.StdoutPipe()

	err := cmd.Start()

	if err != nil {
		return 1, err
	}

	var scanner *bufio.Scanner

	if s.Option.Writer != nil {
		scanner = bufio.NewScanner(io.TeeReader(io.MultiReader(stdout, stderr), s.Option.Writer))
	} else {
		scanner = bufio.NewScanner(io.MultiReader(stdout, stderr))
	}

	for scanner.Scan() {
		m := scanner.Text()

		if s.Option.HandleFunc != nil {
			s.Option.HandleFunc(strings.TrimSpace(m))
		}

		if !s.Option.Quiet {
			fmt.Println(strings.TrimSpace(m))
		}
	}

	var waitStatus syscall.WaitStatus

	if err := cmd.Wait(); err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			// The program has exited with an exit code != 0

			// This works on both Unix and Windows. Although package
			// syscall is generally platform dependent, WaitStatus is
			// defined for both Unix and Windows and in both cases has
			// an ExitStatus() method with the same signature.
			waitStatus = exiterr.Sys().(syscall.WaitStatus)
		}
	}

	waitStatus = cmd.ProcessState.Sys().(syscall.WaitStatus)

	return waitStatus.ExitStatus(), nil
}
