package gocommand

import (
	"context"
	"io"
	"os"
	"os/exec"
	"syscall"
)

// LinuxCommand结构体
type LinuxCommand struct {
}

// LinuxCommand的初始化函数
func NewLinuxCommand() *LinuxCommand {
	return &LinuxCommand{}
}

// 执行命令行并返回结果
// args: 命令行参数
// return: 进程的pid, 命令行结果, 错误消息
func (lc *LinuxCommand) Exec(ctx context.Context, args ...string) (int, string, error) {
	args = append([]string{"-c"}, args...)
	cmd := exec.CommandContext(ctx, os.Getenv("SHELL"), args...)

	cmd.SysProcAttr = &syscall.SysProcAttr{}

	outpip, err := cmd.StdoutPipe()
	if err != nil {
		return 0, "", err
	}
	errpip, err := cmd.StderrPipe()
	if err != nil {
		return 0, "", err
	}
	defer outpip.Close()
	defer errpip.Close()
	err = cmd.Start()
	if err != nil {
		return 0, "", err
	}
	esr, _ := io.ReadAll(errpip)
	out, err := io.ReadAll(outpip)
	if err != nil {
		return 0, "", err
	}

	return cmd.Process.Pid, string(out) + string(esr), nil
}

// 异步执行命令行并通过channel返回结果
// stdout: chan结果
// args: 命令行参数
// return: 进程的pid
// exception: 协程内的命令行发生错误时,会panic异常
func (lc *LinuxCommand) ExecAsync(ctx context.Context, stdout chan string, args ...string) int {
	var pidChan = make(chan int, 1)

	go func() {
		args = append([]string{"-c"}, args...)
		cmd := exec.CommandContext(ctx, os.Getenv("SHELL"), args...)

		cmd.SysProcAttr = &syscall.SysProcAttr{}

		outpip, err := cmd.StdoutPipe()
		if err != nil {
			panic(err)
		}
		defer outpip.Close()
		errpip, err := cmd.StderrPipe()
		if err != nil {
			panic(err)
		}
		defer errpip.Close()

		err = cmd.Start()
		if err != nil {
			panic(err)
		}

		pidChan <- cmd.Process.Pid
		esr, _ := io.ReadAll(errpip)
		out, err := io.ReadAll(outpip)
		if err != nil {
			panic(err)
		}

		stdout <- string(out) + string(esr)
	}()

	return <-pidChan
}

// 执行命令行(忽略返回值)
// args: 命令行参数
// return: 错误消息
func (lc *LinuxCommand) ExecIgnoreResult(ctx context.Context, args ...string) error {
	args = append([]string{"-c"}, args...)
	cmd := exec.CommandContext(ctx, os.Getenv("SHELL"), args...)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{}

	err := cmd.Run()

	return err
}
