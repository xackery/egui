package aseprite

import (
	"encoding/xml"
	"fmt"
	"io"

	"github.com/xackery/egui"
)

// A Reader reads data from an aseprite-data encoded file.
//
// The exported fields can be changed to customize the details before the
// first call to Read or ReadAll.
type Reader struct {
	r io.Reader
}

// NewReader returns a new Reader that reads from r.
func NewReader(r io.Reader) *Reader {
	return &Reader{
		r: r,
	}
}

// ReadAll reads all data from an aseprite file
func (r *Reader) ReadAll() (slices map[string]*egui.Slice, err error) {
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
