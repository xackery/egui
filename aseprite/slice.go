package aseprite

import (
	"encoding/xml"
	"fmt"

	"github.com/xackery/egui"
)

// ReadSlices reads slice data
func (r *Reader) ReadSlices() (slices map[string]*egui.Slice, err error) {
	xr := xml.NewDecoder(r.r)

	type sprite struct {
		Slices []egui.Slice `xml:"slices>slice"`
	}
	sp := &sprite{}
	err = xr.Decode(sp)
	if err != nil {
		return
	}
	slices = make(map[string]*egui.Slice)
	for i := range sp.Slices {
		slice := sp.Slices[i]
		_, ok := slices[slice.Name]
		if ok {
			err = fmt.Errorf("duplicate slice id %s found, must be unique", slice.Name)
			return
		}

		slices[slice.Name] = &slice
	}
	return slices, nil
}
