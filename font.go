package egui

import (
	"fmt"
	"image/color"

	"github.com/golang/freetype/truetype"
	"github.com/hajimehoshi/ebiten"
	"github.com/pkg/errors"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/math/fixed"
	"golang.org/x/text/language"
)

// Font contains related data for loading a font file
type Font struct {
	Face                font.Face
	Height              int
	Name                string
	Language            language.Tag
	boundStringCache    map[font.Face]map[string]*boundStringCacheEntry
	renderingLineHeight int
}

type boundStringCacheEntry struct {
	bounds  *fixed.Rectangle26_6
	advance fixed.Int26_6
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
		Name:                name,
		boundStringCache:    make(map[font.Face]map[string]*boundStringCacheEntry),
		renderingLineHeight: 18,
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

/*
// MeasureSize returns the size of provided text
func (f *Font) MeasureSize(text string) (int, int) {
	w := fixed.I(0)
	h := fixed.I(0)
	for _, l := range strings.Split(strings.TrimRight(text, "\n"), "\n") {
		b, _ := f.boundString(l)
		nw := b.Max.X - b.Min.X
		if nw > w {
			w = nw
		}
		h += fixed.I(f.renderingLineHeight)
	}
	return w.Ceil(), h.Ceil()
}
*/

// DrawText draws a text
func (f *Font) DrawText(dst *ebiten.Image, str string, ox, oy float64, scale float64, textAlign int, color color.Color, displayTextRuneCount int) {
	//f := face(scale, lang)
	/*
		m := f.Metrics()
		oy += (RenderingLineHeight*scale - m.Height.Round()) / 2

		b, _, _ := f.GlyphBounds('.')
		dotX := (-b.Min.X).Floor()

		str = strings.Replace(str, "\r\n", "\n", -1)
		lines := strings.Split(str, "\n")
		linesToShow := strings.Split(string([]rune(str)[:displayTextRuneCount]), "\n")

		for i, l := range linesToShow {
			x := ox + dotX
			y := oy + mplusDotY*scale
			_, a := boundString(f, lines[i])
			switch textAlign {
			case data.TextAlignLeft:
				// do nothing
			case data.TextAlignCenter:
				x -= a.Ceil() / 2
			case data.TextAlignRight:
				x -= a.Ceil()
			default:
				panic(fmt.Sprintf("font: invalid text align: %d", textAlign))
			}

			text.Draw(d, l, f, x, y, color)
			oy += RenderingLineHeight * scale
		}
	*/
}

/*
func (f *Font) boundString(str string) (*fixed.Rectangle26_6, fixed.Int26_6) {
	m, ok := f.boundStringCache[face]
	if !ok {
		m = map[string]*boundStringCacheEntry{}
		f.boundStringCache[face] = m
	}

	entry, ok := m[str]
	if !ok {
		// Delete all entries if the capacity exceeds the limit.
		if len(m) >= 256 {
			for k := range m {
				delete(m, k)
			}
		}

		b, a := font.BoundString(face, str)
		entry = &boundStringCacheEntry{
			bounds:  &b,
			advance: a,
		}
		m[str] = entry
	}

	return entry.bounds, entry.advance
}
*/
