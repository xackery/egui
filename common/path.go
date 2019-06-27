package common

import (
	"container/heap"
	"fmt"
)

var (
	// ErrNoNodesLoaded is returned when a path is called before any nodes are set
	ErrNoNodesLoaded = fmt.Errorf("no nodes are loaded")
	// ErrRouteNotFound is returned when all options are exhausted
	ErrRouteNotFound = fmt.Errorf("route not found")
)

// Path represents a navigation mesh
type Path struct {
	nodes    map[int]map[int]*Node
	isLoaded bool
}

// NewPath returns a new path system
func NewPath() (p *Path) {
	p = &Path{
		nodes: make(map[int]map[int]*Node),
	}
	return
}

// Node returns a node at a coordinate
func (p *Path) Node(x int, y int) *Node {
	if !p.isLoaded {
		return nil
	}
	if x < 0 || y < 0 {
		return nil
	}
	nY, ok := p.nodes[x]
	if !ok {
		return nil
	}
	return nY[y]
}

// NewNode adds a new node
func (p *Path) NewNode(ix int, iy int, isCollider bool, cost float64) {
	x := int(ix)
	y := int(iy)
	p.isLoaded = true
	n := p.Node(int(x), int(y))
	if n != nil {
		n.IsCollider = isCollider
		n.Cost = cost
		return
	}
	_, ok := p.nodes[x]
	if !ok {
		p.nodes[x] = make(map[int]*Node)
	}
	p.nodes[x][y] = &Node{X: int(x), Y: int(y), IsCollider: isCollider, Cost: cost}
	return
}

// Route using astar
func (p *Path) Route(startX, startY, endX, endY int) (path []*Edge, err error) {
	if !p.isLoaded {
		return nil, ErrNoNodesLoaded
	}
	node := p.Node(startX, startY)
	if node == nil {
		err = fmt.Errorf("invalid starting point")
		return
	}
	defer func() {
		if len(path) > 1 {
			path = path[1 : len(path)-1]
			path = append(path, &Edge{Dest: node})
		}
	}()
	var ok bool
	seen := make(map[int]map[int]bool)
	seen[startX] = make(map[int]bool)

	openHeap := make(PriorityQueue, 0)
	heap.Init(&openHeap)
	cameFrom := make(map[int]map[int]*Edge)

	gScore := make(map[int]map[int]float64)
	gScore[startX] = make(map[int]float64)
	fScore := make(map[int]map[int]float64)
	fScore[startX] = make(map[int]float64)

	gScore[startX][startY] = node.Cost
	fScore[startX][startY] = gScore[startX][startY] + node.Heuristic(endX, endY)
	heap.Push(&openHeap, &Item{node: node, priority: 0})

	seen[startX][startY] = true
	var item interface{}
	for {
		item = heap.Pop(&openHeap)
		if item == nil {
			err = fmt.Errorf("out of items")
			return
		}
		node = item.(*Item).node
		if node == nil {
			err = fmt.Errorf("out of items")
			return
		}

		if node.IsCollider {
			continue
		}
		if node.Success(endX, endY) {
			return reconstructPath(cameFrom, node), nil
		}
		for _, edge := range p.Neighbors(node.X, node.Y) {
			if edge == nil {
				continue
			}
			if edge.Dest == nil {
				continue
			}
			adj := edge.Dest
			action := edge.Action
			_, ok = seen[adj.X]
			if ok && seen[adj.X][adj.Y] {
				continue
			}
			if !ok {
				seen[adj.X] = make(map[int]bool)
			}
			seen[adj.X][adj.Y] = true

			_, ok = gScore[adj.X]
			if !ok {
				gScore[adj.X] = make(map[int]float64)
			}
			_, ok = fScore[adj.X]
			if !ok {
				fScore[adj.X] = make(map[int]float64)
			}
			_, ok = cameFrom[adj.X]
			if !ok {
				cameFrom[adj.X] = make(map[int]*Edge)
			}
			// adjacency cost is based on a constant step
			adjCost := node.Cost
			if adjCost == 0 {
				adjCost = 1
			}
			gScore[adj.X][adj.Y] = gScore[node.X][node.Y] + adjCost
			hScore := adj.Heuristic(endX, endY)
			fScore[adj.X][adj.Y] = gScore[adj.X][adj.Y] + hScore
			heap.Push(&openHeap, &Item{node: adj, priority: fScore[adj.X][adj.Y]})
			// reverse the edge for reconstruction
			cameFrom[adj.X][adj.Y] = &Edge{
				Dest:   node,
				Action: action,
				score:  fScore[adj.X][adj.Y],
			}
		}
	}
}

// Neighbors returns edge neighbors of a position
func (p *Path) Neighbors(nodeX int, nodeY int) (neighbors []*Edge) {
	neighbors = []*Edge{
		&Edge{Dest: p.Node(nodeX+1, nodeY), Action: Right},
		&Edge{Dest: p.Node(nodeX-1, nodeY), Action: Left},
		&Edge{Dest: p.Node(nodeX, nodeY+1), Action: Up},
		&Edge{Dest: p.Node(nodeX, nodeY-1), Action: Down},
	}
	return
}

func reconstructPath(cameFrom map[int]map[int]*Edge, node *Node) []*Edge {
	if edge, ok := cameFrom[node.X][node.Y]; ok {
		return append(reconstructPath(cameFrom, edge.Dest), edge)
	}
	return make([]*Edge, 0)
}
