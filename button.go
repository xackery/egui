package egui

import (
	"image/color"
	"time"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/inpututil"
	"github.com/hajimehoshi/ebiten/text"
	"github.com/pkg/errors"
	"github.com/xackery/egui/common"
	"golang.org/x/image/font"
)

// Button represents a UI Button element
type Button struct {
	name               string
	defaultImage       string
	image              *Image
	shape              *common.Rectangle
	text               string
	isEnabled          bool
	isVisible          bool
	alignment          int
	isPressed          bool
	onPressed          func(e *Button)
	onPressFunction    func()
	renderIndex        int64
	isDestroyed        bool
	lerpPosition       *lerpPosition
	lerpColor          *lerpColor
	color              color.Color
	font               *Font
	pressedSliceName   string
	unpressedSliceName string
	ui                 *UI
	isTextShadow       bool
}

// NewButton creates a new button instance
func (u *UI) NewButton(name string, scene string, text string, shape common.Rectangle, textColor color.Color, imageName string, pressedSliceName string, unpressedSliceName string) (*Button, error) {
	if imageName == "" {
		return nil, ErrImageNotFound
	}
	img, err := u.Image(imageName)
	if err != nil {
		return nil, errors.Wrap(err, imageName)
	}

	s, err := u.Scene(scene)
	if err != nil {
		return nil, ErrSceneNotFound
	}

	newShape := common.Rect(shape.Min.X, shape.Min.Y, shape.Max.X, shape.Max.Y)
	e := &Button{
		name:               name,
		image:              img,
		text:               text,
		isEnabled:          true,
		isVisible:          true,
		lerpPosition:       &lerpPosition{},
		lerpColor:          &lerpColor{},
		color:              textColor,
		shape:              &newShape,
		font:               u.defaultFont,
		pressedSliceName:   pressedSliceName,
		unpressedSliceName: unpressedSliceName,
	}

	err = s.AddElement(e)
	if err != nil {
		return nil, err
	}
	return e, nil
}

// Name returns a mob's name
func (e *Button) Name() string {
	return e.name
}

// IsVisible returns true if mob is visible
func (e *Button) IsVisible() bool {
	return e.isVisible
}

// IsEnabled returns true if a button is enabled
func (e *Button) IsEnabled() bool {
	return e.isEnabled
}

// SetEnabled changes if a button is enabled
func (e *Button) SetEnabled(isEnabled bool) {
	e.isEnabled = isEnabled
	return
}

// SetVisible changes the visibility of a button
func (e *Button) SetVisible(isVisible bool) {
	e.isVisible = isVisible
	return
}

// RenderIndex returns the render index of element
func (e *Button) RenderIndex() int64 {
	return e.renderIndex
}

// SetRenderIndex sets the render index of element
func (e *Button) SetRenderIndex(renderIndex int64) {
	e.renderIndex = renderIndex
}

// Update is called during a game update
func (e *Button) update(dt float64) {

	if e.lerpPosition.isEnabled {
		e.shape.Min.X, e.shape.Min.Y = e.lerpPosition.Lerp()
		if !e.lerpPosition.isEnabled {
			if e.lerpPosition.endFunc != nil {
				e.lerpPosition.endFunc()
			}
			if e.lerpPosition.isDestroyed {
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
		if e.shape.Min.X <= fx && fx < e.shape.Max.X && e.shape.Min.Y <= fy && fy < e.shape.Max.Y {
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
		if e.shape.Min.X <= fx && fx < e.shape.Max.X && e.shape.Min.Y <= fy && fy < e.shape.Max.Y {
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
func (e *Button) draw(dst *ebiten.Image) {
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
	op.GeoM.Translate(e.X(), e.Y())
	op.GeoM.Scale(1, 1)

	//opacity := uint8(255)

	if !e.isEnabled {
		op.ColorM.ChangeHSV(0, 0, 1)
		op.ColorM.Scale(0.5, 0.5, 0.5, 1)
	}
	DrawNineSlicing(dst, e.image.ebitenImage, slice.Keys[0], int(e.shape.Dx()), int(e.shape.Dy()), &op.GeoM, &op.ColorM)
	//bounds, _ := font.BoundString(e.font.Face, e.text)
	//w := float64((bounds.Max.X - bounds.Min.X).Ceil())
	//text.Draw(dst, e.text, e.font.Face, int(e.X()), int(e.Y()), e.color)
	bounds, _ := font.BoundString(e.font.Face, e.text)
	w := float64((bounds.Max.X - bounds.Min.X).Ceil())
	x := e.shape.Min.X + (e.shape.Dx()-w)/2
	y := e.shape.Max.Y - (e.shape.Dy()-float64(e.font.Height))/2
	text.Draw(dst, e.text, e.font.Face, int(x), int(y), e.color)

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

	return
}

// SetText changes the text on the button
func (e *Button) SetText(text string) {
	e.text = text
}

// SetOnPressed sets a button state
func (e *Button) SetOnPressed(f func(e *Button)) {
	e.onPressed = f
}

// SetOnPressFunction lets you pass a function without the need of button handling
func (e *Button) SetOnPressFunction(f func()) {
	e.onPressFunction = f
}

// IsDestroyed returns true when the element is flagged for deletion
func (e *Button) IsDestroyed() bool {
	return e.isDestroyed
}

// LerpPosition changes an element's position over duration
func (e *Button) LerpPosition(endPosition common.Vector, duration time.Duration, isDestroyed bool, endFunc func()) {
	e.lerpPosition.start = time.Now()
	e.lerpPosition.startPosition = &common.Vector{X: e.shape.Min.X, Y: e.shape.Min.Y}
	e.lerpPosition.endPosition = &common.Vector{X: endPosition.X, Y: endPosition.Y}
	e.lerpPosition.duration = duration
	e.lerpPosition.isEnabled = true
	e.lerpPosition.endFunc = endFunc
	e.lerpPosition.isDestroyed = isDestroyed
}

// Shape returns an element's X/Y position as well as width/height
func (e *Button) Shape() *common.Rectangle {
	return e.shape
}

// SetShape sets an element's X/Y position as well as width/height
func (e *Button) SetShape(shape common.Rectangle) {
	newShape := common.Rect(shape.Min.X, shape.Min.Y, shape.Max.X, shape.Max.Y)
	e.shape = &newShape
	return
}

// SetIsDestroyed sets an element to be destroyed on next update
func (e *Button) SetIsDestroyed(isDestroyed bool) {
	e.isDestroyed = true
	return
}

// X returns the X position of a button
func (e *Button) X() float64 {
	return e.shape.Min.X
}

// Y returns the Y position of a button
func (e *Button) Y() float64 {
	return e.shape.Min.Y
}
