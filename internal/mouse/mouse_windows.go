//go:build windows

package mouse

import (
	"fmt"
	"syscall"
	"unsafe"
)

type point struct {
	X int32
	Y int32
}

var (
	user32            = syscall.NewLazyDLL("user32.dll")
	procSetCursorPos  = user32.NewProc("SetCursorPos")
	procGetCursorPos  = user32.NewProc("GetCursorPos")
	procGetSystemMtrx = user32.NewProc("GetSystemMetrics")
	procMouseEvent    = user32.NewProc("mouse_event")
)

const (
	smCXScreen = 0
	smCYScreen = 1

	mouseeventfWheel = 0x0800
)

func moveTo(x, y int) error {
	r, _, err := procSetCursorPos.Call(uintptr(x), uintptr(y))
	if r == 0 {
		return fmt.Errorf("SetCursorPos failed: %w", err)
	}
	return nil
}

func scrollVertical(lines int) error {
	if lines == 0 {
		return nil
	}
	const wheelDelta = 120
	data := int32(lines * wheelDelta)
	procMouseEvent.Call(mouseeventfWheel, 0, 0, uintptr(uint32(data)), 0)
	return nil
}

func position() (int, int, error) {
	var p point
	r, _, err := procGetCursorPos.Call(uintptr(unsafe.Pointer(&p)))
	if r == 0 {
		return 0, 0, fmt.Errorf("GetCursorPos failed: %w", err)
	}
	return int(p.X), int(p.Y), nil
}

func screenSize() (int, int, error) {
	w, _, errW := procGetSystemMtrx.Call(smCXScreen)
	h, _, errH := procGetSystemMtrx.Call(smCYScreen)
	if w == 0 || h == 0 {
		return 0, 0, fmt.Errorf("GetSystemMetrics failed: %v / %v", errW, errH)
	}
	return int(w), int(h), nil
}
