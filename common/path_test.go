package common

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGrid(t *testing.T) {
	assert := assert.New(t)
	path := NewPath()
	path.NewNode(0, 0, false, 0)
	path.NewNode(1, 0, true, 0)
	path.NewNode(0, 1, false, 0)
	path.NewNode(1, 1, false, 0)
	path.NewNode(2, 1, false, 0)
	path.NewNode(2, 2, false, 0)
	route, err := path.Path(0, 0, 1, 1)
	assert.NoError(err)
	fmt.Printf("path length %d\n", len(route))
	for _, r := range route {
		fmt.Printf("%s->", r.Action)
	}

	//assert.Empty(route)
	assert.NotEmpty(route)

}
