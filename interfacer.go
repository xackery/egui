package egui

import (
	"time"

	"github.com/hajimehoshi/ebiten"
	"github.com/xackery/egui/common"
)

// Interfacer wraps all user interface elements is a generic user interface wrapper
type Interfacer interface {
	//TargetPositionUpdate(tp *common.Vector, duration time.Duration)
	IsEnabled() bool
	SetEnabled(isEnabled bool)
	IsVisible() bool
	SetVisible(isVisible bool)
	Update(dt float64)
	Draw(screen *ebiten.Image)
	Name() string
	RenderIndex() int64
	SetRenderIndex(renderIndex int64)
	IsDestroyed() bool
	SetIsDestroyed(isDestroyed bool)
	SetText(text string)
	// Shape returns the shape of an element
	Shape() *common.Rectangle
	// SetShape sets the shape of an element
	SetShape(shape common.Rectangle)
	LerpPosition(endPosition common.Vector, duration time.Duration, isDestroyed bool, endFunc func())
}
