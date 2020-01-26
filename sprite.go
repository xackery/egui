package egui

import (
	"image/color"

	"github.com/pkg/errors"
	"github.com/xackery/egui/common"
	"github.com/xackery/egui/element/sprite"
)

// NewSprite creates a new button instance
func (u *UI) NewSprite(name string, scene string, x float64, y float64, imageName string) (*sprite.Element, error) {
	img, err := u.Image(imageName)
	if err != nil {
		return nil, errors.Wrap(err, imageName)
	}

	width, height := img.EbitenImage.Size()

	s, err := u.Scene(scene)
	if err != nil {
		return nil, common.ErrSceneNotFound
	}

	e, err := sprite.New(name, scene, x, y, width, height, color.White, img)
	err = s.AddElement(e)
	if err != nil {
		return nil, err
	}
	return e, nil
}
