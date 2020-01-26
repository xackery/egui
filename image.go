package egui

import (
	"image"
	"io"

	"github.com/hajimehoshi/ebiten"
	"github.com/pkg/errors"
	"github.com/xackery/egui/common"
)

// NewImage adds a new image to egui
func (u *UI) NewImage(name string, f io.Reader, filter ebiten.Filter) (*common.Image, error) {
	img := &common.Image{
		Name:   name,
		Slices: make(map[string]*common.Slice),
	}

	rawImg, _, err := image.Decode(f)
	if err != nil {
		return nil, errors.Wrap(err, "decode")
	}

	img.EbitenImage, err = ebiten.NewImageFromImage(rawImg, ebiten.FilterDefault)
	if err != nil {
		return nil, errors.Wrap(err, "ebiten load")
	}
	err = u.AddImage(img)
	if err != nil {
		return nil, err
	}
	return img, nil
}
