package sprite

import (
	"fmt"
	"image"
	"image/color"
	"time"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/inpututil"
	"github.com/xackery/egui/common"
)

// Element represents a UI Element element. It contains animation data
type Element struct {
	name            string
	image           *common.Image
	x               float64
	y               float64
	width           int
	height          int
	text            string
	isEnabled       bool
	isVisible       bool
	isPressed       bool
	onPressed       func(e *Element)
	onPressFunction func()
	renderIndex     int64
	isDestroyed     bool
	lerpPosition    *common.LerpPosition
	lerpColor       *common.LerpColor
	color           color.Color
	isIdleAnimation bool
	animation       *common.Animation
	scale           float64
	isAnimated      bool
}

// New creates a new element
func New(name string, scene string, x float64, y float64, width int, height int, tintColor color.Color, img *common.Image) (*Element, error) {
	e := &Element{
		name:         name,
		image:        img,
		isEnabled:    true,
		isVisible:    true,
		lerpPosition: new(common.LerpPosition),
		lerpColor:    new(common.LerpColor),
		color:        tintColor,
		x:            x,
		y:            y,
		width:        width,
		height:       height,
		scale:        1,
		isAnimated:   true,
	}
	e.generateAnimation()

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
func (e *Element) Draw(screen *ebiten.Image) {
	if !e.isVisible {
		return
	}
	anim := e.animation

	//opacity := uint8(255)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Reset()
	if !e.isEnabled {
		op.ColorM.ChangeHSV(0, 0, 1)
		op.ColorM.Scale(0.5, 0.5, 0.5, 1)
	}
	var pos []int

	if e.isAnimated {
		if !e.isIdleAnimation {
			anim.Counter++

			ai, ok := anim.Animations[fmt.Sprintf("%d_%s", anim.BundleIndex, anim.CurrentName)]
			if !ok {
				fmt.Println("anim not found")
				//TODO: add a buffer for recent errors
				return
			}
			i := (anim.Counter / anim.Speed) % len(ai)

			animData := ai[i]
			if len(animData) != 2 {
				fmt.Println("animData not found")
				//TODO: add a buffer for recent errors
				return
			}
			pos = anim.Clips[int(animData[0])]
		} else {
			//TODO: idle animation data
		}
	}

	op.GeoM.Translate(e.x, e.y)
	op.GeoM.Scale(e.scale, e.scale)

	if !e.isAnimated {
		screen.DrawImage(e.image.EbitenImage, op)
		return
	}

	if len(pos) != 6 {
		//TODO: add a buffer for recent errors
		return
	}
	r := image.Rect(pos[0], pos[1], pos[2], pos[3])
	op.GeoM.Translate(float64(pos[4]), float64(pos[5]))

	screen.DrawImage(e.image.EbitenImage.SubImage(r).(*ebiten.Image), op)
	e.width = int(e.x + anim.CellWidth + float64(pos[4]))
	e.height = int(e.y + anim.CellHeight + float64(pos[5]))
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
func (e *Element) LerpPosition(endPosition common.Vector, duration time.Duration, isDestroyed bool, endFunc func()) {
	e.lerpPosition.Init(time.Now(), &common.Vector{X: e.x, Y: e.y}, &common.Vector{X: endPosition.X, Y: endPosition.Y}, duration, true, endFunc, isDestroyed)
}

// SetIsDestroyed sets an element to be destroyed on next update
func (e *Element) SetIsDestroyed(isDestroyed bool) {
	e.isDestroyed = true
}

// generateAnimation is called when a new image is loaded.
// It will attempt to split an image into dimensions based on RPGMaker sheets
func (e *Element) generateAnimation() {
	anim := common.Animation{
		Animations: make(map[string][][]float64),
	}

	e.animation = &anim

	bounds := e.image.EbitenImage.Bounds()
	if bounds.Dx() < 26 || bounds.Dy() < 36 {
		if len(e.image.Slices) == 0 {
			e.isAnimated = false
		}
		return
	}
	tileWidth := bounds.Max.X / 3
	tileHeight := bounds.Max.Y / 4
	funkyOffset := 0
	if bounds.Max.X == 312 {
		//funkyOffset = 1
		tileWidth = 26
		tileHeight = 36
	}

	xOffset := 0
	yOffset := 0
	index := 0
	animOrientation := []string{"down", "left", "right", "up"}
	animOrientationIndex := 0

	animName := ""
	for yOffset < bounds.Max.Y {
		for xOffset < bounds.Max.X {
			for y := yOffset; y < yOffset+(tileHeight*4); y += tileHeight {
				for x := xOffset; x < xOffset+(tileWidth*3); x += tileWidth {
					anim.Clips = append(anim.Clips, []int{x, y, x + tileWidth, y + tileHeight, 0, 0})
				}
				animName = fmt.Sprintf("%d_%s", index, animOrientation[animOrientationIndex])
				anim.Animations[animName] = [][]float64{
					[]float64{float64(len(anim.Clips) - (3 - funkyOffset)), 0.2083},
					[]float64{float64(len(anim.Clips) - (2 - funkyOffset)), 0.2083},
					[]float64{float64(len(anim.Clips) - (1 - funkyOffset)), 0.2083},
					[]float64{float64(len(anim.Clips) - (2 - funkyOffset)), 0.2083},
				}

				animOrientationIndex++
			}
			animOrientationIndex = 0
			xOffset += tileWidth * 3
			index++
		}
		xOffset = 0
		yOffset += tileHeight * 4
	}
	anim.BundleCount = index
	anim.CellWidth = float64(tileWidth)
	anim.CellHeight = float64(tileHeight)
	anim.CurrentName = "down"
	anim.Speed = 30
}

// SetAnimation sets animation data
func (e *Element) SetAnimation(anim common.Animation) error {
	e.animation = &common.Animation{
		Counter:     anim.Counter,
		CurrentName: anim.CurrentName,
		Speed:       anim.Speed,
		BundleIndex: anim.BundleIndex,
		CellWidth:   anim.CellWidth,
		CellHeight:  anim.CellHeight,
		Image:       anim.Image,
		Alpha:       anim.Alpha,
		Clips:       anim.Clips,
		BundleCount: anim.BundleCount,
		Animations:  anim.Animations,
	}
	return nil
}

// SetAnimationName sets the current animation group name
func (e *Element) SetAnimationName(name string) {
	e.animation.CurrentName = name
}

// AnimationName returns the current played animation name
func (e *Element) AnimationName() string {
	return e.animation.CurrentName
}

// SetScale sets the scale of a element
func (e *Element) SetScale(scale float64) {
	e.scale = scale
}

// Scale returns the scale of a element. default 1
func (e *Element) Scale() float64 {
	return e.scale
}

// SetIsAnimated flags if the animation information should be honored or not
func (e *Element) SetIsAnimated(isAnimated bool) {
	e.isAnimated = isAnimated
}

// IsAnimated returns true if the element is considered an animation
func (e *Element) IsAnimated() bool {
	return e.isAnimated
}

// SetPosition sets an element's position
func (e *Element) SetPosition(x float64, y float64) {
	e.x = x
	e.y = y
}

// SetBundleIndex sets which version of sprite bundle to render
func (e *Element) SetBundleIndex(index int) error {
	if index > e.animation.BundleCount || index < 0 {
		return fmt.Errorf("index out of range")
	}
	e.animation.BundleIndex = index
	return nil
}

// BundleIndex returns the index sprite rendering
func (e *Element) BundleIndex() int {
	return e.animation.BundleIndex
}
