package mouse

type Point struct {
	X int
	Y int
}

func MoveTo(x, y int) error {
	return moveTo(x, y)
}

func ScrollVertical(lines int) error {
	return scrollVertical(lines)
}

func Position() (Point, error) {
	x, y, err := position()
	if err != nil {
		return Point{}, err
	}
	return Point{X: x, Y: y}, nil
}

func ScreenSize() (Point, error) {
	w, h, err := screenSize()
	if err != nil {
		return Point{}, err
	}
	return Point{X: w, Y: h}, nil
}
