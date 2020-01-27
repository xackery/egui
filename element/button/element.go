package button

import (
	"image/color"
	"time"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/inpututil"
	"github.com/hajimehoshi/ebiten/text"
	"github.com/xackery/egui/common"
)

// Element represents a UI clickable 9slice button.
// A element has text
// A element can register events on click
type Element struct {
	name               string
	image              *common.Image
	x                  float64
	y                  float64
	scale              float64
	width              int
	height             int
	text               string
	isEnabled          bool
	isVisible          bool
	isPressed          bool
	onPressed          func(e *Element)
	onPressFunction    func()
	renderIndex        int64
	isDestroyed        bool
	lerpPosition       *common.LerpPosition
	lerpColor          *common.LerpColor
	color              color.Color
	font               *common.Font
	pressedSliceName   string
	unpressedSliceName string
}

// New creates a new button instance
func New(name string, scene string, text string, x float64, y float64, width int, height int, font *common.Font, textColor color.Color, img *common.Image, pressedSliceName string, unpressedSliceName string) (*Element, error) {

	e := &Element{
		name:               name,
		image:              img,
		text:               text,
		isEnabled:          true,
		isVisible:          true,
		lerpPosition:       new(common.LerpPosition),
		lerpColor:          new(common.LerpColor),
		color:              textColor,
		x:                  x,
		y:                  y,
		width:              width,
		height:             height,
		font:               font,
		pressedSliceName:   pressedSliceName,
		unpressedSliceName: unpressedSliceName,
		scale:              1,
	}

	return e, nil
}

// Name returns a mob's name
func (e *Element) Name() string {
	return e.name
}

// IsVisible returns true if mob is visible
func (e *Element) IsVisible() bool {
	return e.isVisible
}

// IsEnabled returns true if a element is enabled
func (e *Element) IsEnabled() bool {
	return e.isEnabled
}

// SetEnabled changes if a element is enabled
func (e *Element) SetEnabled(isEnabled bool) {
	e.isEnabled = isEnabled
}

// SetVisible changes the visibility of a element
func (e *Element) SetVisible(isVisible bool) {
	e.isVisible = isVisible
}

// RenderIndex returns the render index of element
func (e *Element) RenderIndex() int64 {
	return e.renderIndex
}

// SetRenderIndex sets the render index of element
func (e *Element) SetRenderIndex(renderIndex int64) {
	e.renderIndex = renderIndex
}

// Update is called during a game update
func (e *Element) Update(dt float64) {

	if e.lerpPosition.IsEnabled() {
		e.x, e.y = e.lerpPosition.Lerp()
		if !e.lerpPosition.IsEnabled() {
			if e.lerpPosition.EndFunc() != nil {
				e.lerpPosition.EndFunc()
			}
			if e.lerpPosition.IsDestroyed() {
				e.isDestroyed = true
				return
			}
		}
	}

	isRecentlyPressed := false
	//mobile and desktop use differnet touch devices
	for _, t := range inpututil.JustPressedTouchIDs() {
		x, y := ebiten.TouchPosition(t)
		fx := float64(x)
		fy := float64(y)
		if e.x <= fx && fx < e.x+float64(e.width) && e.y <= fy && fy < e.y+float64(e.height) {
			e.isPressed = true
			isRecentlyPressed = true
		} else {
			e.isPressed = false
		}
	}

	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		fx := float64(x)
		fy := float64(y)
		if e.x <= fx && fx < e.x+float64(e.width) && e.y <= fy && fy < e.y+float64(e.height) {
			e.isPressed = true
			isRecentlyPressed = true
		} else {
			e.isPressed = false
		}
	}

	if !isRecentlyPressed && e.isPressed {
		if e.onPressed != nil {
			e.onPressed(e)
		}
		if e.onPressFunction != nil {
			e.onPressFunction()
		}
		e.isPressed = false
	}
}

// Draw is called during a game update
func (e *Element) Draw(dst *ebiten.Image) {
	if !e.isVisible {
		return
	}
	sliceName := e.unpressedSliceName
	if e.isPressed {
		sliceName = e.pressedSliceName
	}
	slice, err := e.image.Slice(sliceName)
	if err != nil {
		//fmt.Println("slice", sliceName, "not found", err)
		//TODO: handle this error elegantly
		return
	}

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(e.x, e.y)
	op.GeoM.Scale(e.scale, e.scale)

	//opacity := uint8(255)

	if !e.isEnabled {
		op.ColorM.ChangeHSV(0, 0, 1)
		op.ColorM.Scale(0.5, 0.5, 0.5, 1)
	}
	common.DrawNineSlicing(dst, e.image.EbitenImage, slice.Keys[0], e.width, int(e.height), &op.GeoM, &op.ColorM)
	//bounds, _ := font.BoundString(e.font.Face, e.text)
	//w := float64((bounds.Max.X - bounds.Min.X).Ceil())
	//text.Draw(dst, e.text, e.font.Face, int(e.X()), int(e.Y()), e.color)
	//bounds, _ := font.BoundString(e.font.Face, e.text)

	text.Draw(dst, e.text, e.font.Face, int(e.x), int(e.y), e.color)

	/*_, th := e.font.MeasureSize(e.text)
	tx := e.X() * e.ui.tileScale
	tx += e.shape.Dx() * e.ui.tileScale / 2

	ty := e.Y() * e.ui.tileScale
	ty += (e.shape.Dy()*e.ui.tileScale - float64(th)*e.ui.textScale) / 2

	cr, cg, cb, ca := e.color.RGBA()
	r8 := uint8(cr >> 8)
	g8 := uint8(cg >> 8)
	b8 := uint8(cb >> 8)
	a8 := uint8(ca >> 8)
	var c color.Color = color.RGBA{r8, g8, b8, uint8(uint16(a8) * uint16(opacity) / 255)}
	if !e.isEnabled {
		c = color.RGBA{r8, g8, b8, uint8(uint16(a8) * uint16(opacity) / (2 * 255))}
	}
	l := e.font.Language
	if l == language.Und {
		l = e.ui.defaultLanguage
	}
	if e.isTextShadow {
		e.font.DrawText(dst, e.text, tx+e.ui.textScale, ty+e.ui.textScale, e.ui.textScale, 0, color.Black, len([]rune(e.text)))
	}
	e.font.DrawText(dst, e.text, tx, ty, e.ui.textScale, 0, c, len([]rune(e.text)))
	*/

}

// SetText changes the text on the element
func (e *Element) SetText(text string) {
	e.text = text
}

// SetOnPressed sets a element state
func (e *Element) SetOnPressed(f func(e *Element)) {
	e.onPressed = f
}

// SetOnPressFunction lets you pass a function without the need of element handling
func (e *Element) SetOnPressFunction(f func()) {
	e.onPressFunction = f
}

// IsDestroyed returns true when the element is flagged for deletion
func (e *Element) IsDestroyed() bool {
	return e.isDestroyed
}

// LerpPosition changes an element's position over duration
func (e *Element) LerpPosition(endPositionX, endPositionY float64, duration time.Duration, isDestroyed bool, endFunc func()) {
	e.lerpPosition.Init(time.Now(), e.x, e.y, endPositionX, endPositionY, duration, true, endFunc, isDestroyed)
}

// Position returns an element's position
func (e *Element) Position() (float64, float64) {
	return e.x, e.y
}

// SetPosition sets an element's position
func (e *Element) SetPosition(x float64, y float64) {
	e.x = x
	e.y = y
}

// Width returns an element's width
func (e *Element) Width() int {
	return e.width
}

// SetWidth sets an element's width
func (e *Element) SetWidth(width int) {
	e.width = width
}

// Height returns an element's height
func (e *Element) Height() int {
	return e.height
}

// SetHeight sets an element's height
func (e *Element) SetHeight(height int) {
	e.height = height
}

// SetIsDestroyed sets an element to be destroyed on next update
func (e *Element) SetIsDestroyed(isDestroyed bool) {
	e.isDestroyed = true
}
