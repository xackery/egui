package common

import (
	"image/color"
	"time"
)

// LerpColor handles color lerp interpolations
type LerpColor struct {
	start        time.Time
	duration     time.Duration
	startColor   color.Color
	endColor     color.Color
	endFunc      func()
	isDestroyed  bool
	isEnabled    bool
	isEndFuncSet bool
}

// Lerp returns a color
func (lc *LerpColor) Lerp() (newColor color.Color) {
	if !lc.isEnabled {
		newColor = color.White
		return
	}
	if lc.start.Add(lc.duration).Before(time.Now()) {
		lc.isEnabled = false
		newColor = lc.endColor
		return
	}
	elapsed := time.Since(lc.start).Nanoseconds()
	destNano := lc.start.Add(lc.duration).Sub(lc.start).Nanoseconds()
	t := float64(float64(elapsed) / float64(destNano))
	aR, aB, aG, aA := lc.startColor.RGBA()
	bR, bB, bG, bA := lc.endColor.RGBA()
	newColor = color.RGBA64{
		R: uint16((1-t)*float64(aR) + t*float64(bR)),
		G: uint16((1-t)*float64(aG) + t*float64(bG)),
		B: uint16((1-t)*float64(aB) + t*float64(bB)),
		A: uint16((1-t)*float64(aA) + t*float64(bA)),
	}

	return
}

// Init sets up a new lerp
func (lc *LerpColor) Init(start time.Time, startColor color.Color, endColor color.Color, duration time.Duration, endFunc func(), isDestroyedAtEnd bool) {
	lc.start = time.Now()
	lc.startColor = startColor
	lc.endColor = endColor
	lc.duration = duration
	lc.endFunc = endFunc
	lc.isEnabled = true
	lc.isDestroyed = isDestroyedAtEnd
}

// IsEnabled returns if enabled
func (lc *LerpColor) IsEnabled() bool {
	return lc.isEnabled
}

// SetIsEnabled sets if enabled or not
func (lc *LerpColor) SetIsEnabled(isEnabled bool) {
	lc.isEnabled = isEnabled
}

// IsDestroyed returns if enabled
func (lc *LerpColor) IsDestroyed() bool {
	return lc.isDestroyed
}

// SetIsDestroyed sets if destroyed on next frame or not
func (lc *LerpColor) SetIsDestroyed(isDestroyed bool) {
	lc.isDestroyed = isDestroyed
}

// SetEndFunc sets a function to call on end of lerp
func (lc *LerpColor) SetEndFunc(endFunc func()) {
	lc.endFunc = endFunc
	lc.isEndFuncSet = true
}

// EndFunc returns the end function
func (lc *LerpColor) EndFunc() func() {
	return lc.endFunc
}

// IsEndFuncSet returns true if EndFunc exists
func (lc *LerpColor) IsEndFuncSet() bool {
	return lc.isEndFuncSet
}
