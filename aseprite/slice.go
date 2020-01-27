package aseprite

import (
	"encoding/json"
	"fmt"

	"github.com/xackery/egui/common"
)

// ReadSlices reads slice data
func (r *Reader) ReadSlices() (slices map[string]*common.Slice, err error) {
	jr := json.NewDecoder(r.r)

	sp := struct {
		Meta struct {
			Slices []struct {
				Name  string
				Color string
				Keys  []common.SliceKey
			}
		}
	}{}

	err = jr.Decode(&sp)
	if err != nil {
		return
	}
	slices = make(map[string]*common.Slice)
	for _, s := range sp.Meta.Slices {
		slice := new(common.Slice)
		slice.Name = s.Name

		for _, k := range s.Keys {

			/*key := new(common.SliceKey)
			key.CX = float64(k.Center.X)
			key.CX = float64(k.Center.Y)
			key.CH = float64(k.Center.W)
			key.CW = float64(k.Center.H)
			key.Frame = k.Frame
			key.H = float64(k.Bounds.H)
			key.W = float64(k.Bounds.W)
			key.X = float64(k.Bounds.X)
			key.Y = float64(k.Bounds.Y)
			key.PivotX = float64(k.Pivot.X)
			key.PivotY = float64(k.Pivot.Y)
			fmt.Println(key)*/
			slice.Keys = append(slice.Keys, &k)
		}

		_, ok := slices[slice.Name]
		if ok {
			err = fmt.Errorf("duplicate slice id %s found, must be unique", slice.Name)
			return
		}

		slices[slice.Name] = slice
	}
	return slices, nil
}
