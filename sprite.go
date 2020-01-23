package egui

import (
	"fmt"
	"image"
	"image/color"
	"time"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/inpututil"
	"github.com/pkg/errors"
	"github.com/xackery/egui/common"
)

// Sprite represents a UI Sprite element. It contains animation data
type Sprite struct {
	name            string
	image           *Image
	shape           *common.Rectangle
	text            string
	isEnabled       bool
	isVisible       bool
	isPressed       bool
	onPressed       func(e *Sprite)
	onPressFunction func()
	renderIndex     int64
	isDestroyed     bool
	lerpPosition    *lerpPosition
	lerpColor       *lerpColor
	color           color.Color
	isIdleAnimation bool
	animation       *Animation
	scale           float64
	isAnimated      bool
}

// Animation handles animation details
type Animation struct {
	//Counter tracks what animation frame is currently being played
	Counter int
	//Current sprite name being played
	CurrentName string
	//Speed to play animation
	Speed int
	//If the sheet has multiple sprites, which index of the bundle to use
	BundleIndex int
	CellWidth   float64
	CellHeight  float64
	Image       string
	Alpha       string
	Clips       [][]int
	BundleCount int
	Animations  map[string][][]float64
}

// NewSprite creates a new sprite instance
func (u *UI) NewSprite(name string, scene string, shape common.Rectangle, tintColor color.Color, imageName string) (*Sprite, error) {
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
	e := &Sprite{
		name:         name,
		image:        img,
		isEnabled:    true,
		isVisible:    true,
		lerpPosition: &lerpPosition{},
		lerpColor:    &lerpColor{},
		color:        tintColor,
		shape:        &newShape,
		scale:        1,
		isAnimated:   true,
	}
	e.generateAnimation()

	err = s.AddElement(e)
	if err != nil {
		return nil, err
	}
	return e, nil
}

// Name returns a mob's name
func (e *Sprite) Name() string {
	return e.name
}

// IsVisible returns true if mob is visible
func (e *Sprite) IsVisible() bool {
	return e.isVisible
}

// IsEnabled returns true if a sprite is enabled
func (e *Sprite) IsEnabled() bool {
	return e.isEnabled
}

// SetEnabled changes if a sprite is enabled
func (e *Sprite) SetEnabled(isEnabled bool) {
	e.isEnabled = isEnabled
}

// SetVisible changes the visibility of a sprite
func (e *Sprite) SetVisible(isVisible bool) {
	e.isVisible = isVisible
}

// RenderIndex returns the render index of element
func (e *Sprite) RenderIndex() int64 {
	return e.renderIndex
}

// SetRenderIndex sets the render index of element
func (e *Sprite) SetRenderIndex(renderIndex int64) {
	e.renderIndex = renderIndex
}

// Update is called during a game update
func (e *Sprite) update(dt float64) {

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
func (e *Sprite) draw(dst *ebiten.Image) {
	if !e.isVisible {
		return
	}
	anim := e.animation

	op := &ebiten.DrawImageOptions{}
	//opacity := uint8(255)

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

	op.GeoM.Translate(e.X(), e.Y())
	op.GeoM.Scale(e.scale, e.scale)

	if !e.isAnimated {
		dst.DrawImage(e.image.ebitenImage, op)
		return
	}

	if len(pos) != 6 {
		//TODO: add a buffer for recent errors
		return
	}
	r := image.Rect(pos[0], pos[1], pos[2], pos[3])
	op.SourceRect = &r
	op.GeoM.Translate(float64(pos[4]), float64(pos[5]))

	dst.DrawImage(e.image.ebitenImage, op)
	e.shape.Max.X = e.shape.Min.X + anim.CellWidth + float64(pos[4])
	e.shape.Max.Y = e.shape.Min.Y + anim.CellHeight + float64(pos[5])

}

// SetText changes the text on the sprite
func (e *Sprite) SetText(text string) {
	e.text = text
}

// SetOnPressed sets a sprite state
func (e *Sprite) SetOnPressed(f func(e *Sprite)) {
	e.onPressed = f
}

// SetOnPressFunction lets you pass a function without the need of sprite handling
func (e *Sprite) SetOnPressFunction(f func()) {
	e.onPressFunction = f
}

// IsDestroyed returns true when the element is flagged for deletion
func (e *Sprite) IsDestroyed() bool {
	return e.isDestroyed
}

// LerpPosition changes an element's position over duration
func (e *Sprite) LerpPosition(endPosition common.Vector, duration time.Duration, isDestroyed bool, endFunc func()) {
	e.lerpPosition.start = time.Now()
	e.lerpPosition.startPosition = &common.Vector{X: e.shape.Min.X, Y: e.shape.Min.Y}
	e.lerpPosition.endPosition = &common.Vector{X: endPosition.X, Y: endPosition.Y}
	e.lerpPosition.duration = duration
	e.lerpPosition.isEnabled = true
	e.lerpPosition.endFunc = endFunc
	e.lerpPosition.isDestroyed = isDestroyed
}

// Shape returns an element's X/Y position as well as width/height
func (e *Sprite) Shape() *common.Rectangle {
	return e.shape
}

// SetShape sets an element's X/Y position as well as width/height
func (e *Sprite) SetShape(shape common.Rectangle) {
	newShape := common.Rect(shape.Min.X, shape.Min.Y, shape.Max.X, shape.Max.Y)
	e.shape = &newShape
}

// SetIsDestroyed sets an element to be destroyed on next update
func (e *Sprite) SetIsDestroyed(isDestroyed bool) {
	e.isDestroyed = true
}

// X returns the X position of a sprite
func (e *Sprite) X() float64 {
	return e.shape.Min.X
}

// Y returns the Y position of a sprite
func (e *Sprite) Y() float64 {
	return e.shape.Min.Y
}

// generateAnimation is called when a new image is loaded.
// It will attempt to split an image into dimensions based on RPGMaker sheets
func (e *Sprite) generateAnimation() {
	anim := Animation{
		Animations: make(map[string][][]float64),
	}

	e.animation = &anim

	bounds := e.image.ebitenImage.Bounds()
	if bounds.Dx() < 26 || bounds.Dy() < 36 {
		if len(e.image.slices) == 0 {
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
func (e *Sprite) SetAnimation(anim Animation) error {
	e.animation = &Animation{
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
func (e *Sprite) SetAnimationName(name string) {
	e.animation.CurrentName = name
}

// AnimationName returns the current played animation name
func (e *Sprite) AnimationName() string {
	return e.animation.CurrentName
}

// SetScale sets the scale of a sprite
func (e *Sprite) SetScale(scale float64) {
	e.scale = scale
}

// Scale returns the scale of a sprite. default 1
func (e *Sprite) Scale() float64 {
	return e.scale
}

// SetIsAnimated flags if the animation information should be honored or not
func (e *Sprite) SetIsAnimated(isAnimated bool) {
	e.isAnimated = isAnimated
}

// IsAnimated returns true if the sprite is considered an animation
func (e *Sprite) IsAnimated() bool {
	return e.isAnimated
}
