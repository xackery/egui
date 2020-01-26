package common

// Slice represents a 9 slicing instruction set within an image
type Slice struct {
	Name string      `xml:"name"`
	Keys []*SliceKey `xml:"key"`
}

// SliceKey represents each slice's key data
type SliceKey struct {
	Frame  int
	Bounds struct {
		X int
		Y int
		W int
		H int
	}
	Center struct {
		X int
		Y int
		W int
		H int
	}
	Pivot struct {
		X int
		Y int
	}
}
