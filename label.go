package egui

import (
	"image/color"
	"time"

	"github.com/google/uuid"
	"github.com/xackery/egui/common"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/text"
)

// Label represents a UI label element
type Label struct {
	name         string
	shape        *common.Rectangle
	text         string
	isEnabled    bool
	isVisible    bool
	lerpPosition *lerpPosition
	lerpColor    *lerpColor
	color        color.Color

	isDestroyed     bool
	isPressed       bool
	onPressed       func(e *Label)
	onPressFunction func()
	renderIndex     int64
	font            *Font
}

// NewLabel creates a new label instance
func (u *UI) NewLabel(name string, text string, shape common.Rectangle, color color.Color) (e *Label, err error) {

	if name == "" {
		name = uuid.New().String()
	}
	newShape := common.Rect(shape.Min.X, shape.Min.Y, shape.Max.X, shape.Max.Y)
	e = &Label{
		name:         name,
		shape:        &newShape,
		text:         text,
		isEnabled:    true,
		isVisible:    true,
		lerpPosition: &lerpPosition{},
		lerpColor:    &lerpColor{},
		color:        color,
		font:         u.defaultFont,
	}
	err = u.currentScene.AddElement(e)
	if err != nil {
		return
	}
	return
}

// Name returns an element's name
func (e *Label) Name() string {
	return e.name
}

// IsVisible returns true if element should be shown
func (e *Label) IsVisible() bool {
	return e.isVisible
}

// IsEnabled returns true if a button is enabled
func (e *Label) IsEnabled() bool {
	return e.isEnabled
}

// SetEnabled changes if a button is enabled
func (e *Label) SetEnabled(isEnabled bool) {
	e.isEnabled = isEnabled
}

// SetVisible changes the visibility of a button
func (e *Label) SetVisible(isVisible bool) {
	e.isVisible = isVisible
}

// RenderIndex returns the render index of element
func (e *Label) RenderIndex() int64 {
	return e.renderIndex
}

// SetRenderIndex sets the render index of element
func (e *Label) SetRenderIndex(renderIndex int64) {
	e.renderIndex = renderIndex
}

func (e *Label) update(dt float64) {
	if e.lerpColor.enabled {
		e.color = e.lerpColor.Lerp()
		if !e.lerpColor.enabled {
			if e.lerpColor.endFunc != nil {
				e.lerpColor.endFunc()
			}
			if e.lerpColor.isDestroyed {
				e.isDestroyed = true
				return
			}
		}

	}

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

	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		fx := float64(x)
		fy := float64(y)

		if e.shape.Min.X <= fx && fx < e.shape.Max.X && e.shape.Min.Y <= fy && fy < e.shape.Max.Y {
			e.isPressed = true
		} else {
			e.isPressed = false
		}
	} else {
		if e.isPressed {
			if e.onPressed != nil {
				e.onPressed(e)
			}
			if e.onPressFunction != nil {
				e.onPressFunction()
			}
		}
		e.isPressed = false
	}
}

func (e *Label) draw(dst *ebiten.Image) {

	//bounds, _ := font.BoundString(uiInstance.font, e.text)
	//w := float64((bounds.Max.X - bounds.Min.X).Ceil())
	x := e.shape.Min.X //+ (e.shape.Dx()-w)/2
	y := e.shape.Min.Y //- (e.shape.Dy()-float64(uiInstance.fontMHeight))/2
	text.Draw(dst, e.text, e.font.Face, int(x), int(y), e.color)
}

// SetText changes the text on the label
func (e *Label) SetText(text string) {
	e.text = text
}

// SetOnPressed sets a label state
func (e *Label) SetOnPressed(f func(e *Label)) {
	e.onPressed = f
}

// SetOnPressFunction lets you pass a function without the need of label handling
func (e *Label) SetOnPressFunction(f func()) {
	e.onPressFunction = f
}

// LerpColor changes the label's color over duration
func (e *Label) LerpColor(endColor color.Color, duration time.Duration, isDestroyed bool, endFunc func()) {
	e.lerpColor.start = time.Now()
	e.lerpColor.startColor = e.color
	e.lerpColor.endColor = endColor
	e.lerpColor.duration = duration
	e.lerpColor.enabled = true
	e.lerpColor.endFunc = endFunc
	e.lerpColor.isDestroyed = isDestroyed
}

// LerpPosition changes an element's position over duration
func (e *Label) LerpPosition(endPosition common.Vector, duration time.Duration, isDestroyed bool, endFunc func()) {
	e.lerpPosition.start = time.Now()
	e.lerpPosition.startPosition = &common.Vector{X: e.shape.Min.X, Y: e.shape.Min.Y}
	e.lerpPosition.endPosition = &common.Vector{X: endPosition.X, Y: endPosition.Y}
	e.lerpPosition.duration = duration
	e.lerpPosition.isEnabled = true
	e.lerpPosition.endFunc = endFunc
	e.lerpPosition.isDestroyed = isDestroyed
}

// IsDestroyed returns true when the element is flagged for deletion
func (e *Label) IsDestroyed() bool {
	return e.isDestroyed
}

// Shape returns an element's X/Y position as well as width/height
func (e *Label) Shape() *common.Rectangle {
	return e.shape
}

// SetShape sets an element's X/Y position as well as width/height
func (e *Label) SetShape(shape common.Rectangle) {
	newShape := common.Rect(shape.Min.X, shape.Min.Y, shape.Max.X, shape.Max.Y)
	e.shape = &newShape
}

// SetIsDestroyed sets an element to be destroyed on next update
func (e *Label) SetIsDestroyed(isDestroyed bool) {
	e.isDestroyed = true
}
