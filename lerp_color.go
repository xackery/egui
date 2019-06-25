package egui

import (
	"image/color"
	"time"
)

// lerpColor handles color lerp interpolations
type lerpColor struct {
	start       time.Time
	duration    time.Duration
	startColor  color.Color
	endColor    color.Color
	enabled     bool
	endFunc     func()
	isDestroyed bool
}

func (lc *lerpColor) Lerp() (newColor color.Color) {
	if !lc.enabled {
		newColor = color.White
		return
	}
	if lc.start.Add(lc.duration).Before(time.Now()) {
		lc.enabled = false
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
