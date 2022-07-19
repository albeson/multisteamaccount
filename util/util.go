package util

import (
	"log"
	"os"
	"runtime"
	"strings"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

// TH32CS_SNAPPROCESS is described in https://msdn.microsoft.com/de-de/library/windows/desktop/ms682489(v=vs.85).aspx
const TH32CS_SNAPPROCESS = 0x00000002

func init() {

	//explorer := FindProcessByName("explorer.exe")
	//if explorer != nil {
		// found it
		//fmt.Println(explorer.ProcessID)
	//}
}

// WindowsProcess is an implementation of Process for Windows.
type WindowsProcess struct {
	ProcessID       int
	ParentProcessID int
	Exe             string
}

func processes() ([]WindowsProcess, error) {
	handle, err := windows.CreateToolhelp32Snapshot(TH32CS_SNAPPROCESS, 0)
	if err != nil {
		return nil, err
	}
	defer windows.CloseHandle(handle)

	var entry windows.ProcessEntry32
	entry.Size = uint32(unsafe.Sizeof(entry))
	// get the first process
	err = windows.Process32First(handle, &entry)
	if err != nil {
		return nil, err
	}

	results := make([]WindowsProcess, 0, 50)
	for {
		results = append(results, newWindowsProcess(&entry))

		err = windows.Process32Next(handle, &entry)
		if err != nil {
			// windows sends ERROR_NO_MORE_FILES on last process
			if err == syscall.ERROR_NO_MORE_FILES {
				return results, nil
			}
			return nil, err
		}
	}
}

func FindProcessByName(name string) *WindowsProcess {
	processes, err := processes()
	if err != nil {
		log.Fatal(err)
	}
	for _, p := range processes {
		if strings.ToLower(p.Exe) == strings.ToLower(name) {
			return &p
		}
	}
	return nil
}

func newWindowsProcess(e *windows.ProcessEntry32) WindowsProcess {
	// Find when the string ends for decoding
	end := 0
	for {
		if e.ExeFile[end] == 0 {
			break
		}
		end++
	}

	return WindowsProcess{
		ProcessID:       int(e.ProcessID),
		ParentProcessID: int(e.ParentProcessID),
		Exe:             syscall.UTF16ToString(e.ExeFile[:end]),
	}
}

func TerminateProcess(process os.Process, exitCode int) error {
	h, e := syscall.OpenProcess(syscall.PROCESS_TERMINATE, false, uint32(process.Pid))
	if e != nil {
		return os.NewSyscallError("OpenProcess", e)
	}
	defer syscall.CloseHandle(h)
	e = syscall.TerminateProcess(h, uint32(exitCode))
	if e != nil {
		return os.NewSyscallError("TerminateProcess", e)
	}

	runtime.KeepAlive(process)
	return nil
}