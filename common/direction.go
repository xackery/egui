package common

// Direction represents a direction
type Direction int

const (
	// Up is a direction
	Up = Direction(0)
	// Right is a direction
	Right = Direction(1)
	// Down is a direction
	Down = Direction(2)
	// Left is a direction
	Left = Direction(3)
)

func (d Direction) String() string {
	switch d {
	case 0:
		return "up"
	case 1:
		return "right"
	case 2:
		return "down"
	case 3:
		return "left"
	default:
		return "up"
	}
}
