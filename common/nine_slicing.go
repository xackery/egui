package common

import (
	"image"

	"github.com/hajimehoshi/ebiten"
)

// DrawNineSlicing will render slicing data
func DrawNineSlicing(dst, src *ebiten.Image, sliceKey *SliceKey, width int, height int, geoM *ebiten.GeoM, colorM *ebiten.ColorM) {
	partX := int(sliceKey.Center.X)
	partY := int(sliceKey.Center.Y)

	parts := make([]*ebiten.Image, 9)

	for j := 0; j < 3; j++ {
		for i := 0; i < 3; i++ {
			x := i*partX + int(sliceKey.Bounds.X)
			y := j*partY + int(sliceKey.Bounds.Y)
			parts[j*3+i] = src.SubImage(image.Rect(x, y, x+partX, y+partY)).(*ebiten.Image)
		}
	}

	xn, yn := width/partX, height/partY
	op := &ebiten.DrawImageOptions{}
	if colorM != nil {
		op.ColorM.Concat(*colorM)
	}
	for j := 0; j < yn; j++ {
		sy := 0
		switch j {
		case 0:
			sy = 0
		case yn - 1:
			//bottom
			sy = 2
		default:
			//top
			sy = 1
		}
		for i := 0; i < xn; i++ {
			sx := 0
			switch i {
			case 0:
				sx = 0
			case xn - 1:
				//right center
				sx = 2
			default:
				//center
				sx = 1
			}
			op.GeoM.Reset()
			op.GeoM.Translate(float64(i*partX), float64(j*partY))
			op.GeoM.Concat(*geoM)
			dst.DrawImage(parts[sy*3+sx], op)
		}
	}
}
