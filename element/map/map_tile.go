package egui

// MapTile contains the Gid of a specified tile
type MapTile struct {
	GID uint32 `json:"gid,omitempty"`
}

// Index returns a map tile's index
func (mt MapTile) Index() uint32 {
	return uint32(mt.GID &^ tileFlipped)
}

// SetH changes the horizontal flag
func (mt MapTile) SetH(val bool) {
	if val {
		mt.GID |= tileFlippedHorizontal
	} else {
		mt.GID &^= tileFlippedHorizontal
	}
}

// H reads if the tile should be changed horizontal
func (mt MapTile) H() bool {
	return mt.GID&tileFlippedHorizontal != 0
}

// SetV changes the horizontal flag
func (mt MapTile) SetV(val bool) {
	if val {
		mt.GID |= tileFlippedVertical
	} else {
		mt.GID &^= tileFlippedVertical
	}
}

// V returns true if a vertical rotation is in place
func (mt MapTile) V() bool {
	return mt.GID&tileFlippedVertical != 0
}

// SetD changes the horizontal flag
func (mt MapTile) SetD(val bool) {
	if val {
		mt.GID |= tileFlippedDiagonal
	} else {
		mt.GID &^= tileFlippedDiagonal
	}
}

// D returns true if diagonal flagged
func (mt MapTile) D() bool {
	return mt.GID&tileFlippedDiagonal != 0
}
