//+build windows

package scall

import (
	"syscall"
	"unsafe"
)

func SetTerminalTitle(title string) {
	kernel32, err := syscall.LoadLibrary("kernel32.dll")
	if err == nil {
		defer syscall.FreeLibrary(kernel32)
		setConsoleTitle, err := syscall.GetProcAddress(kernel32, "SetConsoleTitleW")
		if err == nil {
			ptr, err := syscall.UTF16PtrFromString(title)
			if err == nil {
				syscall.Syscall(setConsoleTitle, 1, uintptr(unsafe.Pointer(ptr)), 0, 0)
			}
		}
	}
}

func CheckWin10() (ret bool) {
	defer func() {
		if recover() != nil {
			ret = false
		}
	}()
	type OSVERSIONINFOW struct {
		dwOSVersionInfoSize uint32
		dwMajorVersion      uint32
		dwMinorVersion      uint32
		dwBuildNumber       uint32
		dwPlatformId        uint32
		szCSDVersion        [128]byte
	}
	var v OSVERSIONINFOW
	p := uintptr(unsafe.Pointer(&v))
	kernel32 := syscall.NewLazyDLL("ntdll.dll")
	c := kernel32.NewProc("RtlGetVersion")
	r, _, _ := c.Call(p)
	if r == 0 {
		ret = v.dwMajorVersion == 10
	} else {
		ret = false
	}
	return
}
