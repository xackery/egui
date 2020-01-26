package egui

import (
	"image/color"

	"github.com/pkg/errors"
	"github.com/xackery/egui/common"
	"github.com/xackery/egui/element/button"
)

// NewButton creates a new button instance
func (u *UI) NewButton(name string, scene string, text string, x float64, y float64, width int, height int, textColor color.Color, pressedSliceName string, unpressedSliceName string) (*button.Element, error) {
	imageName := "ui"
	img, err := u.Image(imageName)
	if err != nil {
		return nil, errors.Wrap(err, imageName)
	}

	s, err := u.Scene(scene)
	if err != nil {
		return nil, common.ErrSceneNotFound
	}

	e, err := button.New(name, scene, text, x, y, width, height, u.defaultFont, textColor, img, pressedSliceName, unpressedSliceName)
	err = s.AddElement(e)
	if err != nil {
		return nil, err
	}
	return e, nil
}
