package ui

import (
	"fmt"
	"image"
	"sort"

	"github.com/hajimehoshi/ebiten"
)

// Scene represents a layout of ui
type Scene struct {
	isElementsNextUpdateDirty bool
	elementsNextUpdate        elements
	elements                  elements
}

// NewScene initializes a new scene
func NewScene() (s *Scene) {
	s = &Scene{}
	return
}

// Element returns an element based on name
func (s *Scene) Element(name string) (Interfacer, error) {
	if name == "" {
		return nil, ErrElementNameInvalid
	}
	for _, se := range s.elementsNextUpdate {
		if se.NameRead() != name {
			continue
		}
		return se, nil
	}
	return nil, ErrElementNotFound
}

// AddElement adds an element to the scene list
func (s *Scene) AddElement(e Interfacer) error {
	if e.NameRead() == "" {
		return ErrElementNotFound
	}
	for _, se := range s.elementsNextUpdate {
		if se.NameRead() != e.NameRead() {
			continue
		}
		return ErrElementAlreadyExists
	}
	s.elementsNextUpdate = append(s.elementsNextUpdate, e)
	s.isElementsNextUpdateDirty = true
	sort.Sort(elements(s.elementsNextUpdate))
	return nil
}

// RemoveElement flags an element to be removed next update
func (s *Scene) RemoveElement(name string) error {
	if name == "" {
		return ErrElementNameInvalid
	}
	var isFound bool
	for i := range s.elementsNextUpdate {
		if s.elementsNextUpdate[i].NameRead() != name {
			continue
		}
		isFound = true
		s.elementsNextUpdate[i] = s.elementsNextUpdate[len(s.elementsNextUpdate)-1]
		s.elementsNextUpdate = s.elementsNextUpdate[:len(s.elementsNextUpdate)-1]
		break
	}
	if !isFound {
		return ErrElementNotFound
	}
	sort.Sort(elements(s.elementsNextUpdate))
	s.isElementsNextUpdateDirty = true
	return nil
}

func (s *Scene) update(dt float64) {
	if s.isElementsNextUpdateDirty {
		s.elements = s.elementsNextUpdate
		s.isElementsNextUpdateDirty = false
		fmt.Println("updated elements", len(s.elements), "total")
	}

	for _, e := range s.elements {
		e.update(dt)
		if e.IsDestroyed() {
			s.RemoveElement(e.NameRead())
		}
	}
}

func (s *Scene) draw(screen *ebiten.Image) {
	for _, e := range s.elements {
		if !e.IsVisible() {
			continue
		}
		e.draw(screen)
	}
}

func (s *Scene) onResolutionChange(resolution image.Point) {
}
