package ui

import (
	"time"

	"github.com/xackery/egui/common"
)

// lerpPosition handles vector lerp interpolations
type lerpPosition struct {
	start         time.Time
	startPosition *common.Vector
	duration      time.Duration
	endPosition   *common.Vector
	endFunc       func()
	isDestroyed   bool
	isEnabled     bool
}

func (lc *lerpPosition) Lerp() (x float64, y float64) {

	if !lc.isEnabled {
		return lc.endPosition.X, lc.endPosition.Y
	}
	if lc.start.Add(lc.duration).Before(time.Now()) {
		lc.isEnabled = false
		return lc.endPosition.X, lc.endPosition.Y
	}
	elapsed := time.Now().Sub(lc.start).Nanoseconds()
	destNano := lc.start.Add(lc.duration).Sub(lc.start).Nanoseconds()
	t := float64(float64(elapsed) / float64(destNano))
	x = (1-t)*lc.startPosition.X + t*lc.endPosition.X
	y = (1-t)*lc.startPosition.Y + t*lc.endPosition.Y
	return
}
