package common

import "fmt"

// Edge represents an edge of an node
type Edge struct {
	Dest   *Node
	Action Direction
	score  float64
}

func (e *Edge) String() string {
	return fmt.Sprintf("%s (%0.1f) [%s]", e.Action, e.score, e.Dest)
}
