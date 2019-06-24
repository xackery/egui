package egui

import "fmt"

const (
	tileFlippedHorizontal = 0x80000000
	tileFlippedVertical   = 0x40000000
	tileFlippedDiagonal   = 0x20000000
	tileFlipped           = tileFlippedHorizontal | tileFlippedVertical | tileFlippedDiagonal
)

// GID represents a global ID, uses rotation padding on the last 3 bytes
type GID struct {
	gid uint32
}

// NewGID returns a GID
func NewGID(gid uint32) (g GID) {
	g = GID{
		gid: gid,
	}
	return
}

// Index returns an index without flipping data
func (g GID) Index() uint32 {
	return uint32(g.gid &^ tileFlipped)
}

// Value returns the true GID value packed with rotation data
func (g GID) Value() uint32 {
	return g.gid
}

// SetValue sets the true GID value packed with rotation data
func (g GID) SetValue() uint32 {
	return g.gid
}

// SetRotation sets the rotation of a tile clockwise, supports 0, 90, 180, 270
func (g GID) SetRotation(rotation int) (err error) {
	//771        //0 fff
	//2684355437 // 270 ccw tft
	//1073742701 // 180 ftf
	//536871683  // 90 ccw fft
	switch rotation {
	case 0:
		g.SetH(false)
		g.SetV(false)
		g.SetD(false)
	case 90:
		g.SetH(true)
		g.SetV(false)
		g.SetD(true)
	case 180:
		g.SetH(false)
		g.SetV(true)
		g.SetD(false)
	case 270:
		g.SetH(false)
		g.SetV(false)
		g.SetD(true)
	default:
		err = fmt.Errorf("invalid rotation, must be 0, 90, 180, or 270")
	}
	return
}

// Rotation returns the rotation of the object. Can only be 0, 90, 180, or 270
func (g GID) Rotation() (rotation int) {
	if !g.H() && !g.V() && !g.D() {
		return 0
	}
	if g.H() && !g.V() && g.D() {
		return 90
	}
	if !g.H() && g.V() && !g.D() {
		return 180
	}
	if !g.H() && !g.V() && g.D() {
		return 270
	}
	return 0
}

// SetH changes the horizontal flag
func (g GID) SetH(val bool) {
	if val {
		g.gid |= tileFlippedHorizontal
	} else {
		g.gid &^= tileFlippedHorizontal
	}
}

// H reads if the tile should be changed horizontal
func (g GID) H() bool {
	return g.gid&tileFlippedHorizontal != 0
}

// SetV changes the horizontal flag
func (g GID) SetV(val bool) {
	if val {
		g.gid |= tileFlippedVertical
	} else {
		g.gid &^= tileFlippedVertical
	}
}

// V returns true if a vertical rotation is in place
func (g GID) V() bool {
	return g.gid&tileFlippedVertical != 0
}

// SetD changes the horizontal flag
func (g GID) SetD(val bool) {
	if val {
		g.gid |= tileFlippedDiagonal
	} else {
		g.gid &^= tileFlippedDiagonal
	}
}

// D returns true if diagonal flagged
func (g GID) D() bool {
	return g.gid&tileFlippedDiagonal != 0
}

// Index strips a gid of rotation data
func Index(gid uint32) uint32 {
	return uint32(gid &^ tileFlipped)
}
