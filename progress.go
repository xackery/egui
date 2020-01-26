package egui

import (
	"image/color"

	"github.com/pkg/errors"
	"github.com/xackery/egui/common"
	"github.com/xackery/egui/element/progress"
)

// NewProgress creates a new progress bar instance
func (u *UI) NewProgress(name string, scene string, text string, x float64, y float64, width int, height int, progressImageName string, fillImageName string) (*progress.Element, error) {

	imageName := "ui"
	img, err := u.Image(imageName)
	if err != nil {
		return nil, errors.Wrap(err, imageName)
	}

	s, err := u.Scene(scene)
	if err != nil {
		return nil, common.ErrSceneNotFound
	}

	e, err := progress.New(name, scene, text, x, y, width, height, u.defaultFont, color.White, img, progressImageName, fillImageName)
	err = s.AddElement(e)
	if err != nil {
		return nil, err
	}
	return e, nil
}
