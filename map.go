package egui

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"time"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/inpututil"
	"github.com/pkg/errors"
	"github.com/xackery/egui/common"
)

// Map represents a UI Map element
type Map struct {
	name            string
	defaultImage    string
	image           *Image
	shape           *common.Rectangle
	text            string
	isEnabled       bool
	isVisible       bool
	alignment       int
	isPressed       bool
	onPressed       func(e *Map)
	onPressFunction func()
	renderIndex     int64
	isDestroyed     bool
	lerpPosition    *lerpPosition
	lerpColor       *lerpColor
	color           color.Color
	ui              *UI
	isTextShadow    bool
	isIdleAnimation bool
	data            *MapData
}

// MapData contains map related information
type MapData struct {
	Source          string `json:"source,omitempty"`
	Width           int64  `json:"width,omitempty"`
	Height          int64  `json:"height,omitempty"`
	TileWidth       int64  `json:"tile_width,omitempty"`
	TileHeight      int64  `json:"tile_height,omitempty"`
	TileCount       int64  `json:"tile_count,omitempty"`
	TileSheetWidth  int64
	TileSheetHeight int64
	TileFrames      []image.Rectangle
	Layers          []MapLayer    `json:"layers,omitempty"`
	Colliders       []MapCollider `json:"colliders,omitempty"`
}

// MapLayer contains layer data
type MapLayer struct {
	Name    string    `json:"name,omitempty"`
	Opacity float32   `json:"opacity,omitempty"`
	Tiles   []MapTile `json:"tiles,omitempty"`
}

// MapCollider is a boolean if a collider is true or not
type MapCollider struct {
	IsCollider bool    `json:"is_collider,omitempty"`
	Cost       float32 `json:"cost,omitempty"`
}

// NewMap creates a new map instance
func (u *UI) NewMap(name string, scene string, shape common.Rectangle, tintColor color.Color, imageName string) (*Map, error) {
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
	e := &Map{
		name:         name,
		image:        img,
		isEnabled:    true,
		isVisible:    true,
		lerpPosition: &lerpPosition{},
		lerpColor:    &lerpColor{},
		color:        tintColor,
		shape:        &newShape,
		data:         &MapData{},
	}

	err = s.AddElement(e)
	if err != nil {
		return nil, err
	}
	return e, nil
}

// Name returns a mob's name
func (e *Map) Name() string {
	return e.name
}

// IsVisible returns true if mob is visible
func (e *Map) IsVisible() bool {
	return e.isVisible
}

// IsEnabled returns true if a map is enabled
func (e *Map) IsEnabled() bool {
	return e.isEnabled
}

// SetEnabled changes if a map is enabled
func (e *Map) SetEnabled(isEnabled bool) {
	e.isEnabled = isEnabled
	return
}

// SetVisible changes the visibility of a map
func (e *Map) SetVisible(isVisible bool) {
	e.isVisible = isVisible
	return
}

// RenderIndex returns the render index of element
func (e *Map) RenderIndex() int64 {
	return e.renderIndex
}

// SetRenderIndex sets the render index of element
func (e *Map) SetRenderIndex(renderIndex int64) {
	e.renderIndex = renderIndex
}

// Update is called during a game update
func (e *Map) update(dt float64) {

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
func (e *Map) draw(dst *ebiten.Image) {
	if !e.isVisible {
		return
	}

	var x, y float64

	//for i := len(e.data.Layers) - 1; i >= 0; i-- {
	for i := 0; i < len(e.data.Layers); i++ {
		x = 0
		y = 0
		l := e.data.Layers[i]
		/*if filter == "fg" && !strings.Contains(l.Name, "fg") {
			continue
		}*/
		for j := 0; j < len(l.Tiles); j++ {
			t := l.Tiles[j]
			if t.Index() == 0 {
				x += float64(e.data.TileWidth)
				if x >= float64(e.data.TileSheetWidth) {
					x = 0
					y += float64(e.data.TileHeight)
				}
				continue
			}

			op := &ebiten.DrawImageOptions{}
			h := t.H()
			v := t.V()
			d := t.D()
			if h {
				switch {
				case v && d:
					op.GeoM.Rotate(math.Pi / 2)
					op.GeoM.Translate(8, 0)
					op.GeoM.Scale(1.0, -1.0)
					op.GeoM.Translate(0, 8)
				case v && !d:
					op.GeoM.Scale(-1.0, -1.0)
					op.GeoM.Translate(8, 8)
				case !v && d:
					op.GeoM.Rotate(math.Pi / 2)
					op.GeoM.Translate(8, 0)
				case !v && !d:
					op.GeoM.Scale(-1.0, 1.0)
					op.GeoM.Translate(8, 0)
				}
			} else if v {
				switch {
				case !h && d:
					op.GeoM.Rotate(math.Pi / 2)
					op.GeoM.Translate(8, 0)
					op.GeoM.Scale(-1.0, -1.0)
					op.GeoM.Translate(8, 8)
				case !h && !d:
					op.GeoM.Scale(1.0, -1.0)
					op.GeoM.Translate(0, 8)
				}
			} else if d {
				switch {
				case !v && !h:
					op.GeoM.Rotate(math.Pi / 2)
					op.GeoM.Translate(8, 0)
					op.GeoM.Scale(-1.0, 1.0)
					op.GeoM.Translate(8, 0)
				}
			}

			op.GeoM.Translate(float64(x), float64(y))
			op.GeoM.Translate(e.X(), e.Y())
			if len(e.data.TileFrames) < int(t.Index()) {
				fmt.Println("index out of range:", t.Index(), ">", len(e.data.TileFrames))
			}
			r := e.data.TileFrames[t.Index()]
			op.SourceRect = &r
			dst.DrawImage(e.image.ebitenImage, op)

			x += float64(e.data.TileWidth)
			if x >= float64(e.data.TileSheetWidth) {
				x = 0
				y += float64(e.data.TileHeight)
			}
		}
	}

	return
}

// SetText changes the text on the map
func (e *Map) SetText(text string) {
	e.text = text
}

// SetOnPressed sets a map state
func (e *Map) SetOnPressed(f func(e *Map)) {
	e.onPressed = f
}

// SetOnPressFunction lets you pass a function without the need of map handling
func (e *Map) SetOnPressFunction(f func()) {
	e.onPressFunction = f
}

// IsDestroyed returns true when the element is flagged for deletion
func (e *Map) IsDestroyed() bool {
	return e.isDestroyed
}

// LerpPosition changes an element's position over duration
func (e *Map) LerpPosition(endPosition common.Vector, duration time.Duration, isDestroyed bool, endFunc func()) {
	e.lerpPosition.start = time.Now()
	e.lerpPosition.startPosition = &common.Vector{X: e.shape.Min.X, Y: e.shape.Min.Y}
	e.lerpPosition.endPosition = &common.Vector{X: endPosition.X, Y: endPosition.Y}
	e.lerpPosition.duration = duration
	e.lerpPosition.isEnabled = true
	e.lerpPosition.endFunc = endFunc
	e.lerpPosition.isDestroyed = isDestroyed
}

// Shape returns an element's X/Y position as well as width/height
func (e *Map) Shape() *common.Rectangle {
	return e.shape
}

// SetShape sets an element's X/Y position as well as width/height
func (e *Map) SetShape(shape common.Rectangle) {
	newShape := common.Rect(shape.Min.X, shape.Min.Y, shape.Max.X, shape.Max.Y)
	e.shape = &newShape
	return
}

// SetIsDestroyed sets an element to be destroyed on next update
func (e *Map) SetIsDestroyed(isDestroyed bool) {
	e.isDestroyed = true
	return
}

// X returns the X position of a map
func (e *Map) X() float64 {
	return e.shape.Min.X
}

// Y returns the Y position of a map
func (e *Map) Y() float64 {
	return e.shape.Min.Y
}

// SetData sets a map's data
func (e *Map) SetData(data MapData) error {
	e.data = &data
	return nil
}
