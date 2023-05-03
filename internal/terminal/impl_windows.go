package terminal

import (
	"fmt"
	"github.com/mattn/go-isatty"
	"golang.org/x/sys/windows"
	"os"
	"syscall"
)

func configureOutputHandle(f *os.File) (*OutputStream, error) {
	ret := &OutputStream{
		File: f,
	}

	if fd := f.Fd(); isatty.IsTerminal(fd) {

		err := SetConsoleOutputCP(CP_UTF8)
		if err != nil {
			return nil, fmt.Errorf("failed to set the console to UTF-8 mode; you may need to use a newer version of Windows: %s", err)
		}

		ret.getColumns = getColumnsWindowsConsole
		var mode uint32
		err = windows.GetConsoleMode(windows.Handle(fd), &mode)
		if err != nil {
			return ret, nil // We'll treat this as success but without VT support
		}
		mode |= windows.ENABLE_VIRTUAL_TERMINAL_PROCESSING
		err = windows.SetConsoleMode(windows.Handle(fd), mode)
		if err != nil {
			return ret, nil // We'll treat this as success but without VT support
		}

		ret.isTerminal = staticTrue
		return ret, nil
	} else if isatty.IsCygwinTerminal(fd) {

		ret.isTerminal = staticTrue
		// TODO: Is it possible to detect the width of these fake terminals?
		return ret, nil
	}

	return ret, nil
}

func configureInputHandle(f *os.File) (*InputStream, error) {
	ret := &InputStream{
		File: f,
	}

	if fd := f.Fd(); isatty.IsTerminal(fd) {
		err := SetConsoleCP(CP_UTF8)
		if err != nil {
			return nil, fmt.Errorf("failed to set the console to UTF-8 mode; you may need to use a newer version of Windows: %s", err)
		}
		ret.isTerminal = staticTrue
		return ret, nil
	} else if isatty.IsCygwinTerminal(fd) {
		ret.isTerminal = staticTrue
		return ret, nil
	}

	return ret, nil
}

func getColumnsWindowsConsole(f *os.File) int {
	var info windows.ConsoleScreenBufferInfo
	err := windows.GetConsoleScreenBufferInfo(windows.Handle(f.Fd()), &info)
	if err != nil {
		return defaultColumns
	}
	return int(info.Size.X)
}

var modkernel32 = windows.NewLazySystemDLL("kernel32.dll")
var procSetConsoleCP = modkernel32.NewProc("SetConsoleCP")
var procSetConsoleOutputCP = modkernel32.NewProc("SetConsoleOutputCP")

const CP_UTF8 = 65001

func SetConsoleCP(codepageID uint32) (err error) {
	r1, _, e1 := syscall.Syscall(procSetConsoleCP.Addr(), 1, uintptr(codepageID), 0, 0)
	if r1 == 0 {
		err = e1
	}
	return
}

func SetConsoleOutputCP(codepageID uint32) (err error) {
	r1, _, e1 := syscall.Syscall(procSetConsoleOutputCP.Addr(), 1, uintptr(codepageID), 0, 0)
	if r1 == 0 {
		err = e1
	}
	return
}

func staticTrue(f *os.File) bool {
	return true
}
