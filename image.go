package egui

import (
	"fmt"
	"image"
	"io"

	"github.com/hajimehoshi/ebiten"
	"github.com/pkg/errors"
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
	name        string
	ebitenImage *ebiten.Image
	slices      map[string]*Slice
	animation   *Animation
}

// NewImage adds a new image to egui
func (u *UI) NewImage(name string, f io.Reader, filter ebiten.Filter) (*Image, error) {
	img := &Image{
		name:   name,
		slices: make(map[string]*Slice),
	}

	rawImg, _, err := image.Decode(f)
	if err != nil {
		return nil, errors.Wrap(err, "decode")
	}

	img.ebitenImage, err = ebiten.NewImageFromImage(rawImg, ebiten.FilterDefault)
	if err != nil {
		return nil, errors.Wrap(err, "ebiten load")
	}

	err = u.AddImage(img)
	if err != nil {
		return nil, err
	}
	return img, nil
}

// Slice returns a 9 slice based on name
func (img *Image) Slice(name string) (*Slice, error) {
	if name == "" {
		return nil, ErrSliceNameInvalid
	}
	s, ok := img.slices[name]
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
	_, ok := img.slices[s.Name]
	if ok {
		return ErrSliceAlreadyExists
	}
	img.slices[s.Name] = s
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
	_, ok := img.slices[s.Name]
	if !ok {
		return ErrSliceNotFound
	}
	delete(img.slices, s.Name)
	return nil
}

// Name returns the image name
func (img *Image) Name() string {
	return img.name
}

// EbitenImage returns the ebiten image
func (img *Image) EbitenImage() *ebiten.Image {
	return img.ebitenImage
}

// Slices return image slice info
func (img *Image) Slices() map[string]*Slice {
	return img.slices
}
