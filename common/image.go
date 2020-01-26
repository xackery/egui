package common

import (
	"fmt"

	"github.com/hajimehoshi/ebiten"
)

var (
	// ErrSliceNameInvalid is returned when a slice name has invalid characters or too short
	ErrSliceNameInvalid = fmt.Errorf("slice name invalid")
	// ErrSliceAlreadyExists is returned when a slice already exists
	ErrSliceAlreadyExists = fmt.Errorf("slice already exists")
	// ErrSliceNotFound is returned when a slice is not loaded into the UI
	ErrSliceNotFound = fmt.Errorf("slice not found")
	// ErrFontNameInvalid is returned when a font name has invalid characters or too short
)

// Image is a base type representing a nebiten image with added details on how to render.
// An image is not rendered, instead placed into a cache
type Image struct {
	Name        string
	EbitenImage *ebiten.Image
	Slices      map[string]*Slice
	Animation   *Animation
}

// Slice returns a 9 slice based on name
func (img *Image) Slice(name string) (*Slice, error) {
	if name == "" {
		return nil, ErrSliceNameInvalid
	}
	s, ok := img.Slices[name]
	if !ok {
		return nil, ErrSliceNotFound
	}
	return s, nil
}

// AddSlice adds a 9slicing to an image
func (img *Image) AddSlice(s *Slice) error {
	if s == nil {
		return ErrSliceNameInvalid
	}
	if s.Name == "" {
		return ErrSliceNameInvalid
	}
	_, ok := img.Slices[s.Name]
	if ok {
		return ErrSliceAlreadyExists
	}
	img.Slices[s.Name] = s
	return nil
}

// RemoveSlice removes 9slicing data from an image
func (img *Image) RemoveSlice(s *Slice) error {
	if s == nil {
		return ErrSliceNameInvalid
	}
	if s.Name == "" {
		return ErrSliceNameInvalid
	}
	_, ok := img.Slices[s.Name]
	if !ok {
		return ErrSliceNotFound
	}
	delete(img.Slices, s.Name)
	return nil
}
