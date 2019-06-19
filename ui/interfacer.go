package ui

import (
	"time"

	"github.com/hajimehoshi/ebiten"
	"github.com/xackery/egui/common"
)

// Interfacer is a generic user interface wrapper
type Interfacer interface {
	//TargetPositionUpdate(tp *common.Vector, duration time.Duration)
	EnabledRead() bool
	EnabledUpdate(isEnabled bool)
	VisibleRead() bool
	VisibleUpdate(isVisible bool)
	update(dt float64)
	IsVisible() bool
	draw(screen *ebiten.Image)
	NameRead() string
	RenderIndexRead() int64
	RenderIndexUpdate(renderIndex int64)
	IsDestroyed() bool
	// ShapeRead returns the shape of an element
	ShapeRead() *common.Rectangle
	// ShapeUpdate sets the shape of an element
	ShapeUpdate(shape *common.Rectangle)
	LerpPosition(endPosition *common.Vector, duration time.Duration, isDestroyed bool, endFunc func())
	TextUpdate(string)
}
