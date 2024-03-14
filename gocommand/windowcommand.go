package gocommand

import (
	"context"
	"io"
	"os"
	"os/exec"
	"syscall"

	"golang.org/x/text/encoding/simplifiedchinese"
)

// WindowsCommand结构体
type WindowsCommand struct {
}

// WindowsCommand的初始化函数
func NewWindowsCommand() *WindowsCommand {
	return &WindowsCommand{}
}

type Charset string

const (
	UTF8    = Charset("UTF-8")
	GB18030 = Charset("GB18030")
)

func ConvertByte2String(byte []byte, charset Charset) string {
	var str string
	switch charset {
	case GB18030:
		var decodeBytes, _ = simplifiedchinese.GB18030.NewDecoder().Bytes(byte)
		str = string(decodeBytes)
	case UTF8:
		fallthrough
	default:
		str = string(byte)
	}
	return str
}

// 执行命令行并返回结果
// args: 命令行参数
// return: 进程的pid, 命令行结果, 错误消息
func (lc *WindowsCommand) Exec(ctx context.Context, args ...string) (int, string, error) {
	args = append([]string{"-c"}, args...)
	cmd := exec.CommandContext(ctx, "cmd", args...)

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
	byte2String := ConvertByte2String(out, "GB18030")
	byte2String_err := ConvertByte2String(esr, "GB18030")
	return cmd.Process.Pid, byte2String + byte2String_err, nil
}

// 异步执行命令行并通过channel返回结果
// stdout: chan结果
// args: 命令行参数
// return: 进程的pid
// exception: 协程内的命令行发生错误时,会panic异常
func (lc *WindowsCommand) ExecAsync(ctx context.Context, stdout chan string, args ...string) int {
	var pidChan = make(chan int, 1)

	go func() {
		args = append([]string{"-c"}, args...)
		cmd := exec.CommandContext(ctx, "cmd", args...)

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
		byte2String := ConvertByte2String(out, "GB18030")
		byte2String_err := ConvertByte2String(esr, "GB18030")
		stdout <- byte2String + byte2String_err
	}()

	return <-pidChan
}

// 执行命令行(忽略返回值)
// args: 命令行参数
// return: 错误消息
func (lc *WindowsCommand) ExecIgnoreResult(ctx context.Context, args ...string) error {
	args = append([]string{"-c"}, args...)
	cmd := exec.CommandContext(ctx, "cmd", args...)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{}

	err := cmd.Run()

	return err
}
