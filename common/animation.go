package common

// Animation handles animation details
type Animation struct {
	//Counter tracks what animation frame is currently being played
	Counter int
	//Current sprite name being played
	CurrentName string
	//Speed to play animation
	Speed int
	//If the sheet has multiple sprites, which index of the bundle to use
	BundleIndex int
	CellWidth   float64
	CellHeight  float64
	Image       string
	Alpha       string
	Clips       [][]int
	BundleCount int
	Animations  map[string][][]float64
}
