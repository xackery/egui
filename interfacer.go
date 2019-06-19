package egui

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
	// Shape returns the shape of an element
	Shape() *common.Rectangle
	// SetShape sets the shape of an element
	SetShape(shape common.Rectangle)
	LerpPosition(endPosition common.Vector, duration time.Duration, isDestroyed bool, endFunc func())
	TextUpdate(string)
}
