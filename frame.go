package egui

import (
	"image"
	"time"

	"github.com/xackery/egui/common"

	"github.com/hajimehoshi/ebiten"
)

// Frame represents a UI Frame element
type Frame struct {
	name            string
	image           *ebiten.Image
	shape           *common.Rectangle
	text            string
	drawRect        image.Rectangle
	isEnabled       bool
	isVisible       bool
	alignment       int
	isPressed       bool
	onPressed       func(e *Frame)
	onPressFunction func()
	renderIndex     int64
	isDestroyed     bool

	lerpPosition *lerpPosition
	lerpColor    *lerpColor
}

// NewFrame creates a new Frame instance
func (u *UI) NewFrame(name string, imageName string, text string, shape *common.Rectangle) (e *Frame, err error) {
	img, err := u.Image(imageName)
	if err != nil {
		return
	}

	e = &Frame{
		name:         name,
		image:        img,
		drawRect:     image.Rect(32, 16, 48, 32),
		shape:        shape,
		text:         text,
		isEnabled:    true,
		isVisible:    true,
		lerpPosition: &lerpPosition{},
		lerpColor:    &lerpColor{},
	}
	return
}

// NameRead returns a mob's name
func (e *Frame) NameRead() string {
	return e.name
}

// IsVisible returns true if mob is visible
func (e *Frame) IsVisible() bool {
	return e.VisibleRead()
}

// RenderIndexRead returns the render index of element
func (e *Frame) RenderIndexRead() int64 {
	return e.renderIndex
}

// RenderIndexUpdate sets the render index of element
func (e *Frame) RenderIndexUpdate(renderIndex int64) {
	e.renderIndex = renderIndex
}

// AlignmentUpdate changes the alignment of the element
func (e *Frame) AlignmentUpdate(alignment int) {
	e.alignment = alignment
}

// AlignmentRead returns the alignment style
func (e *Frame) AlignmentRead() int {
	return e.alignment
}

// EnabledRead returns true if a Frame is enabled
func (e *Frame) EnabledRead() bool {
	return e.isEnabled
}

// EnabledUpdate changes if a Frame is enabled
func (e *Frame) EnabledUpdate(isEnabled bool) {
	e.isEnabled = isEnabled
	return
}

// VisibleRead returns true if a Frame is visible
func (e *Frame) VisibleRead() bool {
	return e.isVisible
}

// VisibleUpdate changes the visibility of a Frame
func (e *Frame) VisibleUpdate(isVisible bool) {
	e.isVisible = isVisible
	return
}

// Update is called during a game update
func (e *Frame) update(dt float64) {

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

}

// Draw is called during a game update
func (e *Frame) draw(dst *ebiten.Image) {

	dstRect := &common.Rectangle{}

	e.drawNinePatch(dst, dstRect, e.shape)

	/*	bounds, _ := font.BoundString(uiInstance.font, e.text)
		w := float64((bounds.Max.X - bounds.Min.X).Ceil())
		x := e.shape.Min.X + (e.shape.Dx()-w)/2
		y := e.shape.Max.Y - (e.shape.Dy()-float64(uiInstance.fontMHeight))/2
		text.Draw(dst, e.text, uiInstance.font, int(x), int(y), color.White)*/
}

// TextUpdate changes the text on the Frame
func (e *Frame) TextUpdate(text string) {
	e.text = text
}

// SetOnPressed sets a Frame state
func (e *Frame) SetOnPressed(f func(e *Frame)) {
	e.onPressed = f
}

// SetOnPressFunction lets you pass a function without the need of Frame handling
func (e *Frame) SetOnPressFunction(f func()) {
	e.onPressFunction = f
}

func (e *Frame) drawNinePatches(dst *ebiten.Image, dstRect image.Rectangle, srcRect image.Rectangle) {
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
			dx := 0
			dy := 0
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
			r := image.Rect(sx, sy, sx+sw, sy+sh)
			op.SourceRect = &r
			op.GeoM.Translate(e.shape.Min.X, e.shape.Min.Y)

			dst.DrawImage(e.image, op)
		}
	}
}

// IsDestroyed returns true when the element is flagged for deletion
func (e *Frame) IsDestroyed() bool {
	return e.isDestroyed
}

// LerpPosition changes an element's position over duration
func (e *Frame) LerpPosition(endPosition common.Vector, duration time.Duration, isDestroyed bool, endFunc func()) {
	e.lerpPosition.start = time.Now()
	e.lerpPosition.startPosition = &common.Vector{X: e.shape.Min.X, Y: e.shape.Min.Y}
	e.lerpPosition.endPosition = &common.Vector{X: endPosition.X, Y: endPosition.Y}
	e.lerpPosition.duration = duration
	e.lerpPosition.isEnabled = true
	e.lerpPosition.endFunc = endFunc
	e.lerpPosition.isDestroyed = isDestroyed
}

func (e *Frame) drawNinePatch(dst *ebiten.Image, dstRect *common.Rectangle, srcRect *common.Rectangle) {
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

// Shape returns an element's X/Y position as well as width/height
func (e *Frame) Shape() *common.Rectangle {
	return e.shape
}

// SetShape sets an element's X/Y position as well as width/height
func (e *Frame) SetShape(shape common.Rectangle) {
	newShape := common.Rect(shape.Min.X, shape.Min.Y, shape.Max.X, shape.Max.Y)
	e.shape = &newShape
	return
}
