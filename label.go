package egui

import (
	"image/color"

	"github.com/pkg/errors"
	"github.com/xackery/egui/common"
	"github.com/xackery/egui/element/label"
)

// NewLabel creates a new label instance
func (u *UI) NewLabel(name string, scene string, text string, x float64, y float64, textColor color.Color) (*label.Element, error) {

	img, err := u.Image("ui")
	if err != nil {
		return nil, errors.Wrap(err, "ui")
	}

	s, err := u.Scene(scene)
	if err != nil {
		return nil, common.ErrSceneNotFound
	}

	e, err := label.New(name, scene, text, x, y, u.defaultFont, textColor, img)
	err = s.AddElement(e)
	if err != nil {
		return nil, err
	}
	return e, nil
}
