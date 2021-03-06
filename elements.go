package egui

import (
	"github.com/xackery/egui/element"
)

type elements []element.Interfacer

// Len is part of sort.Interface.
func (e elements) Len() int {
	return len(e)
}

// Swap is part of sort.Interface.
func (e elements) Swap(i, j int) {
	e[i], e[j] = e[j], e[i]
}

// Less is part of sort.Interface. We use count as the value to sort by
func (e elements) Less(i, j int) bool {
	return e[i].RenderIndex() < e[j].RenderIndex()
}
