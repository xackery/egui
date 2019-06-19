package egui

import (
	"image"
	"image/color"
	"time"

	"golang.org/x/image/font"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/inpututil"
	"github.com/hajimehoshi/ebiten/text"
	"github.com/xackery/egui/common"
)

// Button represents a UI Button element
type Button struct {
	name            string
	defaultResource string
	image           *ebiten.Image
	shape           *common.Rectangle
	text            string
	pressedRect     image.Rectangle
	unpressedRect   image.Rectangle
	isEnabled       bool
	isVisible       bool
	alignment       int
	isPressed       bool
	onPressed       func(e *Button)
	onPressFunction func()
	renderIndex     int64
	isDestroyed     bool
	lerpPosition    *lerpPosition
	lerpColor       *lerpColor
	color           color.Color
	//inherited by ui by default
	font font.Face
	//inherited by ui by default
	fontMHeight int
}

// NewButton creates a new button instance
func (u *UI) NewButton(name string, resourceName string, scene string, text string, shape *common.Rectangle, textColor color.Color) (*Button, error) {
	if resourceName == "" {
		return nil, ErrResourceNotFound
	}
	img, err := u.Resource(resourceName)
	if err != nil {
		return nil, err
	}

	s, err := u.Scene(scene)
	if err != nil {
		return nil, ErrSceneNotFound
	}

	e := &Button{
		name:          name,
		image:         img,
		pressedRect:   image.Rect(16, 0, 32, 16),
		unpressedRect: image.Rect(0, 0, 16, 16),
		text:          text,
		isEnabled:     true,
		isVisible:     true,
		lerpPosition:  &lerpPosition{},
		lerpColor:     &lerpColor{},
		color:         textColor,
		shape:         shape,
		font:          u.font,
	}

	err = s.AddElement(e)
	if err != nil {
		return nil, err
	}
	return e, nil
}

// NameRead returns a mob's name
func (e *Button) NameRead() string {
	return e.name
}

// IsVisible returns true if mob is visible
func (e *Button) IsVisible() bool {
	return e.VisibleRead()
}

// AlignmentUpdate changes the alignment of the element
func (e *Button) AlignmentUpdate(alignment int) {
	e.alignment = alignment
}

// AlignmentRead returns the alignment style
func (e *Button) AlignmentRead() int {
	return e.alignment
}

// EnabledRead returns true if a button is enabled
func (e *Button) EnabledRead() bool {
	return e.isEnabled
}

// EnabledUpdate changes if a button is enabled
func (e *Button) EnabledUpdate(isEnabled bool) {
	e.isEnabled = isEnabled
	return
}

// VisibleRead returns true if a button is visible
func (e *Button) VisibleRead() bool {
	return e.isVisible
}

// VisibleUpdate changes the visibility of a button
func (e *Button) VisibleUpdate(isVisible bool) {
	e.isVisible = isVisible
	return
}

// RenderIndexRead returns the render index of element
func (e *Button) RenderIndexRead() int64 {
	return e.renderIndex
}

// RenderIndexUpdate sets the render index of element
func (e *Button) RenderIndexUpdate(renderIndex int64) {
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

	srcRect := e.unpressedRect
	if e.isPressed {
		srcRect = e.pressedRect
	}

	e.drawNinePatch(dst, e.shape, common.RectImageCopy(srcRect))

	bounds, _ := font.BoundString(e.font, e.text)
	w := float64((bounds.Max.X - bounds.Min.X).Ceil())
	x := e.shape.Min.X + (e.shape.Dx()-w)/2
	y := e.shape.Max.Y - (e.shape.Dy()-float64(e.fontMHeight))/2
	text.Draw(dst, e.text, e.font, int(x), int(y), e.color)
}

// TextUpdate changes the text on the button
func (e *Button) TextUpdate(text string) {
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
func (e *Button) LerpPosition(endPosition *common.Vector, duration time.Duration, isDestroyed bool, endFunc func()) {
	e.lerpPosition.start = time.Now()
	e.lerpPosition.startPosition = &common.Vector{X: e.shape.Min.X, Y: e.shape.Min.Y}
	e.lerpPosition.endPosition = endPosition
	e.lerpPosition.duration = duration
	e.lerpPosition.isEnabled = true
	e.lerpPosition.endFunc = endFunc
	e.lerpPosition.isDestroyed = isDestroyed
}

func (e *Button) drawNinePatch(dst *ebiten.Image, dstRect *common.Rectangle, srcRect *common.Rectangle) {

	srcX := srcRect.Min.X
	srcY := srcRect.Min.Y
	srcW := srcRect.Dx()
	srcH := srcRect.Dy()

	dstX := dstRect.Min.X
	dstY := dstRect.Min.Y
	dstW := dstRect.Dx()
	dstH := dstRect.Dy()

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

			dst.DrawImage(e.image, op)
		}
	}
}

// ShapeRead returns an element's X/Y position as well as width/height
func (e *Button) ShapeRead() *common.Rectangle {
	return e.shape
}

// ShapeUpdate sets an element's X/Y position as well as width/height
func (e *Button) ShapeUpdate(shape *common.Rectangle) {
	e.shape = shape
	return
}
