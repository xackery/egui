package common

import (
	"time"
)

// LerpPosition handles vector lerp interpolations
type LerpPosition struct {
	start          time.Time
	startPositionX float64
	startPositionY float64
	duration       time.Duration
	endPositionX   float64
	endPositionY   float64
	endFunc        func()
	isEndFuncSet   bool
	isDestroyed    bool
	isEnabled      bool
}

// Lerp returns a position
func (lc *LerpPosition) Lerp() (x float64, y float64) {

	if !lc.isEnabled {
		return lc.endPositionX, lc.endPositionY
	}
	if lc.start.Add(lc.duration).Before(time.Now()) {
		lc.isEnabled = false
		return lc.endPositionX, lc.endPositionY
	}

	elapsed := time.Since(lc.start).Nanoseconds()
	destNano := lc.start.Add(lc.duration).Sub(lc.start).Nanoseconds()
	t := float64(float64(elapsed) / float64(destNano))
	x = (1-t)*lc.startPositionX + t*lc.endPositionX
	y = (1-t)*lc.startPositionY + t*lc.endPositionY
	return
}

// IsEnabled returns if enabled
func (lc *LerpPosition) IsEnabled() bool {
	return lc.isEnabled
}

// SetIsEnabled sets if enabled or not
func (lc *LerpPosition) SetIsEnabled(isEnabled bool) {
	lc.isEnabled = isEnabled
}

// IsDestroyed returns if enabled
func (lc *LerpPosition) IsDestroyed() bool {
	return lc.isDestroyed
}

// SetIsDestroyed sets if destroyed on next frame or not
func (lc *LerpPosition) SetIsDestroyed(isDestroyed bool) {
	lc.isDestroyed = isDestroyed
}

// SetEndFunc sets a function to call on end of lerp
func (lc *LerpPosition) SetEndFunc(endFunc func()) {
	lc.endFunc = endFunc
	lc.isEndFuncSet = true
}

// EndFunc returns the end function
func (lc *LerpPosition) EndFunc() func() {
	return lc.endFunc
}

// IsEndFuncSet returns true if EndFunc exists
func (lc *LerpPosition) IsEndFuncSet() bool {
	return lc.isEndFuncSet
}

// Init sets up a new lerp
func (lc *LerpPosition) Init(start time.Time, startPositionX, startPositionY, endPositionX, endPositionY float64, duration time.Duration, isEnabled bool, endFunc func(), isDestroyedAtEnd bool) {
	lc.start = time.Now()
	lc.startPositionX = startPositionX
	lc.startPositionY = startPositionY
	lc.endPositionX = endPositionX
	lc.endPositionY = endPositionY
	lc.duration = duration
	lc.isEnabled = true
	lc.endFunc = endFunc
	lc.isDestroyed = isDestroyedAtEnd
}
