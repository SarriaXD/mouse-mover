//go:build darwin

package mouse

/*
#cgo LDFLAGS: -framework ApplicationServices
#include <ApplicationServices/ApplicationServices.h>

static CGEventRef createScrollEvent(int lines) {
	return CGEventCreateScrollWheelEvent(NULL, kCGScrollEventUnitLine, 1, lines);
}
*/
import "C"
import "fmt"

func moveTo(x, y int) error {
	event := C.CGEventCreateMouseEvent(C.CGEventSourceRef(0), C.kCGEventMouseMoved, C.CGPoint{C.double(x), C.double(y)}, C.kCGMouseButtonLeft)
	if event == 0 {
		return fmt.Errorf("CGEventCreateMouseEvent returned nil")
	}
	defer C.CFRelease(C.CFTypeRef(event))
	C.CGEventPost(C.kCGHIDEventTap, event)
	return nil
}

func scrollVertical(lines int) error {
	event := C.createScrollEvent(C.int(lines))
	if event == 0 {
		return fmt.Errorf("CGEventCreateScrollWheelEvent returned nil")
	}
	defer C.CFRelease(C.CFTypeRef(event))
	C.CGEventPost(C.kCGHIDEventTap, event)
	return nil
}

func position() (int, int, error) {
	event := C.CGEventCreate(C.CGEventSourceRef(0))
	if event == 0 {
		return 0, 0, fmt.Errorf("CGEventCreate returned nil")
	}
	defer C.CFRelease(C.CFTypeRef(event))
	loc := C.CGEventGetLocation(event)
	return int(loc.x), int(loc.y), nil
}

func screenSize() (int, int, error) {
	bounds := C.CGDisplayBounds(C.CGMainDisplayID())
	return int(bounds.size.width), int(bounds.size.height), nil
}
