package aseprite

import (
	"io"
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
