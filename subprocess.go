package goutils

import (
	"bufio"
	"context"
	"fmt"
	logger "github.com/sirupsen/logrus"
	"io"
	"os/exec"
	"strings"
	"syscall"
	"time"
)

type ProcessOutLineHandle func(line string)

type SubProcessOption struct {
	Debug      bool
	Quiet      bool
	Timeout    time.Duration
	HandleFunc ProcessOutLineHandle
	ShellExec  string
}

type SubProcess struct {
	Option   SubProcessOption
	Commands []string
}

func NewSubProcess() *SubProcess {
	return &SubProcess{
		Option: SubProcessOption{
			Debug:   false,
			Quiet:   false,
			Timeout: 60 * time.Second,
		},
	}
}

func NewSubProcessWithOption(option SubProcessOption) *SubProcess {
	return &SubProcess{
		Option: option,
	}
}

func (s *SubProcess) String() string {
	return fmt.Sprintf("Timeout: %v, Quiet: %v, Commands: %s", s.Option.Timeout, s.Option.Quiet, strings.Join(s.Commands, " && "))
}

func (s *SubProcess) AddCommand(commands ...string) *SubProcess {
	for _, v := range commands {
		s.Commands = append(s.Commands, v)
	}

	return s
}

func (s *SubProcess) ClearCommand() {
	s.Commands = nil
}

func (s *SubProcess) PrintCommands() *SubProcess {
	fmt.Println(fmt.Sprintf("\x1b[1;33m[COMMAND]\x1b[0m %s", strings.Join(s.Commands, " && ")))

	return s
}

// Run 执行系统命令。
func (s *SubProcess) Run() (int, error) {
	return s.RunWithWriter(nil)
}

// RunWithWriter 执行系统命令并同时输出到 io.Writer。
func (s *SubProcess) RunWithWriter(w io.Writer) (int, error) {
	if s.Option.Debug {
		s.PrintCommands()
	}

	var cmd *exec.Cmd

	execBin := "/bin/sh"

	if s.Option.ShellExec != "" {
		execBin = s.Option.ShellExec
	}

	if s.Option.Timeout > 0 {
		ctx, cancel := context.WithTimeout(context.Background(), s.Option.Timeout)
		defer cancel()
		cmd = exec.CommandContext(ctx, execBin, "-c", strings.Join(s.Commands, " && "))
	} else {
		cmd = exec.Command(execBin, "-c", strings.Join(s.Commands, " && "))
	}

	stderr, _ := cmd.StderrPipe()
	stdout, _ := cmd.StdoutPipe()

	err := cmd.Start()

	if err != nil {
		return 1, err
	}

	var scanner *bufio.Scanner

	if w != nil {
		scanner = bufio.NewScanner(io.TeeReader(io.MultiReader(stdout, stderr), w))
	} else {
		scanner = bufio.NewScanner(io.MultiReader(stdout, stderr))
	}
	// scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		m := scanner.Text()

		if s.Option.HandleFunc != nil {
			s.Option.HandleFunc(strings.TrimSpace(m))
		}

		if !s.Option.Quiet {
			logger.Info(m)
		}
	}

	if err := cmd.Wait(); err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			// The program has exited with an exit code != 0

			// This works on both Unix and Windows. Although package
			// syscall is generally platform dependent, WaitStatus is
			// defined for both Unix and Windows and in both cases has
			// an ExitStatus() method with the same signature.
			status, ok := exiterr.Sys().(syscall.WaitStatus)

			if ok {
				// logger.Printf("Exit Status: %d", status.ExitStatus())
				return status.ExitStatus(), exiterr
			} else {
				return status.ExitStatus(), err
			}
		} else {
			return 2, err
		}
	}

	return 0, nil
}
