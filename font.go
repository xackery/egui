package egui

import (
	"fmt"

	"github.com/golang/freetype/truetype"
	"github.com/pkg/errors"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
)

// Font contains related data for loading a font file
type Font struct {
	Face   font.Face
	Height int
	Name   string
}

// NewFontTTF instantiates a truetype font. truetype.Parse() can be used to load a TTF to fontData
func (u *UI) NewFontTTF(name string, fontData []byte, opts *truetype.Options, r rune) (*Font, error) {
	if opts == nil {
		opts = &truetype.Options{Size: 12, DPI: 72, Hinting: font.HintingFull}
	}

	tt, err := truetype.Parse(goregular.TTF)
	if err != nil {
		return nil, errors.Wrap(err, "parse ttf font")
	}
	f := &Font{
		Name: name,
	}
	f.Face = truetype.NewFace(tt, opts)
	b, _, ok := f.Face.GlyphBounds(r)
	if !ok {
		return nil, fmt.Errorf("calibrate glyph bounds failed")
	}
	f.Height = (b.Max.Y - b.Min.Y).Ceil()

	err = u.AddFont(f)
	if err != nil {
		return nil, err
	}
	return f, nil
}
