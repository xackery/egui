package main

import (
	"fmt"
	"image"
	"os"

	"github.com/golang/freetype/truetype"
	"github.com/pkg/errors"
	"github.com/xackery/egui"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
)

func main() {
	err := run()
	if err != nil {
		fmt.Println("failed to run", err.Error())
		os.Exit(1)
	}
}

func run() error {
	tt, err := truetype.Parse(goregular.TTF)
	if err != nil {
		return errors.Wrap(err, "parse ttf font")
	}
	uiFont := truetype.NewFace(tt, &truetype.Options{
		Size:    12,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	b, _, ok := uiFont.GlyphBounds('M')
	if !ok {
		return fmt.Errorf("calibrate glyph bounds")
	}
	uiFontMHeight := (b.Max.Y - b.Min.Y).Ceil()
	ui, err := egui.NewUI(uiFont, uiFontMHeight, image.Point{X: 640, Y: 480})
	if err != nil {
		return errors.Wrap(err, "start ui")
	}
	fmt.Println(ui)
	return nil
}
