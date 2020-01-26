package egui

import (
	"image"
	"sort"

	"github.com/hajimehoshi/ebiten"
	"github.com/xackery/egui/common"
	"github.com/xackery/egui/element"
)

// Scene represents a layout of ui
type Scene struct {
	isElementsNextUpdateDirty bool
	elementsNextUpdate        elements
	elements                  elements
}

// NewScene initializes a new scene
func (ui *UI) NewScene(name string) (*Scene, error) {
	s := &Scene{}
	err := ui.AddScene(name, s)
	if err != nil {
		return nil, err
	}
	return s, nil
}

// Element returns an element based on name
func (s *Scene) Element(name string) (element.Interfacer, error) {
	if name == "" {
		return nil, common.ErrElementNameInvalid
	}
	for _, se := range s.elementsNextUpdate {
		if se.Name() != name {
			continue
		}
		return se, nil
	}
	return nil, common.ErrElementNotFound
}

// AddElement adds an element to the scene list
func (s *Scene) AddElement(e element.Interfacer) error {
	if e.Name() == "" {
		return common.ErrElementNotFound
	}
	for _, se := range s.elementsNextUpdate {
		if se.Name() != e.Name() {
			continue
		}
		return common.ErrElementAlreadyExists
	}
	s.elementsNextUpdate = append(s.elementsNextUpdate, e)
	s.isElementsNextUpdateDirty = true
	sort.Sort(elements(s.elementsNextUpdate))
	return nil
}

// RemoveElement flags an element to be removed next update
func (s *Scene) RemoveElement(name string) error {
	if name == "" {
		return common.ErrElementNameInvalid
	}
	var isFound bool
	for i := range s.elementsNextUpdate {
		if s.elementsNextUpdate[i].Name() != name {
			continue
		}
		isFound = true
		s.elementsNextUpdate[i] = s.elementsNextUpdate[len(s.elementsNextUpdate)-1]
		s.elementsNextUpdate = s.elementsNextUpdate[:len(s.elementsNextUpdate)-1]
		break
	}
	if !isFound {
		return common.ErrElementNotFound
	}
	sort.Sort(elements(s.elementsNextUpdate))
	s.isElementsNextUpdateDirty = true
	return nil
}

// Update is called during a frame update
func (s *Scene) Update(dt float64) {
	if s.isElementsNextUpdateDirty {
		s.elements = s.elementsNextUpdate
		s.isElementsNextUpdateDirty = false
	}

	for _, e := range s.elements {
		e.Update(dt)
		if e.IsDestroyed() {
			s.RemoveElement(e.Name())
		}
	}
}

// Draw renders on a destination image
func (s *Scene) Draw(screen *ebiten.Image) {
	for _, e := range s.elements {
		if !e.IsVisible() {
			continue
		}
		e.Draw(screen)
	}
}

func (s *Scene) onResolutionChange(resolution image.Point) {
}
