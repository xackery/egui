package egui

import (
	"image"
	"image/color"
	"time"

	"golang.org/x/image/font"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/inpututil"
	"github.com/hajimehoshi/ebiten/text"
	"github.com/pkg/errors"
	"github.com/xackery/egui/common"
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

	e.drawNinePatch(dst)

	bounds, _ := font.BoundString(e.font.Face, e.text)
	w := float64((bounds.Max.X - bounds.Min.X).Ceil())
	x := e.shape.Min.X + (e.shape.Dx()-w)/2
	y := e.shape.Max.Y - (e.shape.Dy()-float64(e.font.Height))/2
	text.Draw(dst, e.text, e.font.Face, int(x), int(y), e.color)
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

func (e *Button) drawNinePatch(dst *ebiten.Image) {

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

	srcX := slice.Keys[0].X
	srcY := slice.Keys[0].Y
	srcW := slice.Keys[0].W
	srcH := slice.Keys[0].H

	dstX := e.shape.Min.X
	dstY := e.shape.Min.Y
	dstW := e.shape.Dx()
	dstH := e.shape.Dy()

	op := &ebiten.DrawImageOptions{}
	for j := 0; j < 3; j++ {
		for i := 0; i < 3; i++ {
			op.GeoM.Reset()
			sx := srcX
			sy := srcY
			sw := srcW / 4
			sh := srcH / 4
			dx := float64(0)
			dy := float64(0)
			dw := sw
			dh := sh
			switch i {
			case 1:
				sx = srcX + srcW/3
				sw = srcW / 3
				dx = srcW / 4
				dw = dstW - 2*srcW/4
			case 2:
				sx = srcX + 3*srcW/4
				dx = dstW - srcW/4
			}
			switch j {
			case 1:
				sy = srcY + srcH/4
				sh = srcH / 2
				dy = srcH / 4
				dh = dstH - 2*srcH/4
			case 2:
				sy = srcY + 3*srcH/4
				dy = dstH - srcH/4
			}

			op.GeoM.Scale(float64(dw)/float64(sw), float64(dh)/float64(sh))

			op.GeoM.Translate(float64(dx), float64(dy))
			op.GeoM.Translate(float64(dstX), float64(dstY))
			r := image.Rect(int(sx), int(sy), int(sx+sw), int(sy+sh))
			op.SourceRect = &r
			op.GeoM.Translate(e.shape.Min.X, e.shape.Min.Y)

			dst.DrawImage(e.image.ebitenImage, op)
		}
	}
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
