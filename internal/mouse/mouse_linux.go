//go:build linux

package mouse

import "fmt"

func moveTo(x, y int) error {
	return fmt.Errorf("linux support is not implemented yet (planned: X11)")
}

func scrollVertical(lines int) error {
	return fmt.Errorf("linux support is not implemented yet (planned: X11)")
}

func position() (int, int, error) {
	return 0, 0, fmt.Errorf("linux support is not implemented yet (planned: X11)")
}

func screenSize() (int, int, error) {
	return 0, 0, fmt.Errorf("linux support is not implemented yet (planned: X11)")
}
