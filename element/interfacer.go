package element

import (
	"time"

	"github.com/hajimehoshi/ebiten"
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
	LerpPosition(endPositionX, endpositionY float64, duration time.Duration, isDestroyed bool, endFunc func())
}
