package common

import (
	"fmt"
	"image"
)

// Rectangle is a 2D rectangle aligned with the axes of the coordinate system. It is defined by two
// points, Min and Max.
//
// The invariant should hold, that Max's components are greater or equal than Min's components
// respectively.
type Rectangle struct {
	Min, Max Vector
}

// Rect returns a rectangle
func Rect(x0 float64, y0 float64, x1 float64, y1 float64) Rectangle {
	return Rectangle{Min: Vector{X: x0, Y: y0}, Max: Vector{X: x1, Y: y1}}
}

// String returns the string representation of the rectangle
func (r *Rectangle) String() string {
	return fmt.Sprintf("Rect(%.0f, %.0f, %.0f, %.0f)", r.Min.X, r.Min.Y, r.Max.X, r.Max.Y)
}

// Dx returns r's width
func (r *Rectangle) Dx() float64 {
	return r.Max.X - r.Min.X
}

// Dy returns r's height
func (r *Rectangle) Dy() float64 {
	return r.Max.Y - r.Min.Y
}

// RectImageCopy converts an image rectangle to model rect
func RectImageCopy(srcRect image.Rectangle) *Rectangle {
	return &Rectangle{Min: Vector{X: float64(srcRect.Min.X), Y: float64(srcRect.Min.Y)}, Max: Vector{X: float64(srcRect.Max.X), Y: float64(srcRect.Max.Y)}}
}

/*
// R returns a new Rect with given the Min and Max coordinates.
//
// Note that the returned rectangle is not automatically normalized.
func R(minX, minY, maxX, maxY float64) Rect {
	return Rect{
		Min: Vector{minX, minY},
		Max: Vector{maxX, maxY},
	}
}

// String returns the string representation of the Rect.
//
//   r := pixel.R(100, 50, 200, 300)
//   r.String()     // returns "Rect(100, 50, 200, 300)"
//   fmt.Println(r) // Rect(100, 50, 200, 300)
func (r Rect) String() string {
	return fmt.Sprintf("Rect(%v, %v, %v, %v)", r.Min.X, r.Min.Y, r.Max.X, r.Max.Y)
}

// Norm returns the Rect in normal form, such that Max is component-wise greater or equal than Min.
func (r Rect) Norm() Rect {
	return Rect{
		Min: Vector{
			math.Min(r.Min.X, r.Max.X),
			math.Min(r.Min.Y, r.Max.Y),
		},
		Max: Vector{
			math.Max(r.Min.X, r.Max.X),
			math.Max(r.Min.Y, r.Max.Y),
		},
	}
}

// W returns the width of the Rect.
func (r Rect) W() float64 {
	return r.Max.X - r.Min.X
}

// H returns the height of the Rect.
func (r Rect) H() float64 {
	return r.Max.Y - r.Min.Y
}

// Size returns the vector of width and height of the Rect.
func (r Rect) Size() *Vector {
	return V(r.W(), r.H())
}

// Area returns the area of r. If r is not normalized, area may be negative.
func (r Rect) Area() float64 {
	return r.W() * r.H()
}

// Center returns the position of the center of the Rect.
func (r Rect) Center() *Vector {
	return Lerp(r.Min, r.Max, 0.5)
}

// Moved returns the Rect moved (both Min and Max) by the given vector delta.
func (r Rect) Moved(delta Vector) Rect {
	return Rect{
		Min: r.Min.Add(delta),
		Max: r.Max.Add(delta),
	}
}

// Resized returns the Rect resized to the given size while keeping the position of the given
// anchor.
//
//   r.Resized(r.Min, size)      // resizes while keeping the position of the lower-left corner
//   r.Resized(r.Max, size)      // same with the top-right corner
//   r.Resized(r.Center(), size) // resizes around the center
//
// This function does not make sense for resizing a rectangle of zero area and will panic. Use
// ResizedMin in the case of zero area.
func (r Rect) Resized(anchor, size Vector) Rect {
	if r.W()*r.H() == 0 {
		panic(fmt.Errorf("(%T).Resize: zero area", r))
	}
	fraction := Vector{size.X / r.W(), size.Y / r.H()}
	return Rect{
		Min: anchor.Add(r.Min.Sub(anchor).ScaledXY(fraction)),
		Max: anchor.Add(r.Max.Sub(anchor).ScaledXY(fraction)),
	}
}

// ResizedMin returns the Rect resized to the given size while keeping the position of the Rect's
// Min.
//
// Sizes of zero area are safe here.
func (r Rect) ResizedMin(size Vector) Rect {
	return Rect{
		Min: r.Min,
		Max: r.Min.Add(size),
	}
}

// Contains checks whether a vector u is contained within this Rect (including it's borders).
func (r Rect) Contains(v *Vector) bool {
	return r.Min.X <= nv.X && nv.X <= r.Max.X && r.Min.Y <= nv.Y && nv.Y <= r.Max.Y
}

// Union returns the minimal Rect which covers both r and s. Rects r and s must be normalized.
func (r Rect) Union(s Rect) Rect {
	return R(
		math.Min(r.Min.X, s.Min.X),
		math.Min(r.Min.Y, s.Min.Y),
		math.Max(r.Max.X, s.Max.X),
		math.Max(r.Max.Y, s.Max.Y),
	)
}

// Intersect returns the maximal Rect which is covered by both r and s. Rects r and s must be normalized.
//
// If r and s don't overlap, this function returns R(0, 0, 0, 0).
func (r Rect) Intersect(s Rect) Rect {
	t := R(
		math.Max(r.Min.X, s.Min.X),
		math.Max(r.Min.Y, s.Min.Y),
		math.Min(r.Max.X, s.Max.X),
		math.Min(r.Max.Y, s.Max.Y),
	)
	if t.Min.X >= t.Max.X || t.Min.Y >= t.Max.Y {
		return Rect{}
	}
	return t
}

// Matrix is a 2x3 affine matrix that can be used for all kinds of spatial transforms, such
// as movement, scaling and rotations.
//
// Matrix has a handful of useful methods, each of which adds a transformation to the matrix. For
// example:
//
//   pixel.IM.Moved(pixel.V(100, 200)).Rotated(pixel.ZV, math.Pi/2)
//
// This code creates a Matrix that first moves everything by 100 units horizontally and 200 units
// vertically and then rotates everything by 90 degrees around the origin.
//
// Layout is:
// [0] [2] [4]
// [1] [3] [5]
//  0   0   1  (implicit row)
type Matrix [6]float64

// IM stands for identity matrix. Does nothing, no transformation.
var IM = Matrix{1, 0, 0, 1, 0, 0}

// String returns a string representation of the Matrix.
//
//   m := pixel.IM
//   fmt.Println(m) // Matrix(1 0 0 | 0 1 0)
func (m Matrix) String() string {
	return fmt.Sprintf(
		"Matrix(%v %v %v | %v %v %v)",
		m[0], m[2], m[4],
		m[1], m[3], m[5],
	)
}

// Moved moves everything by the delta vector.
func (m Matrix) Moved(delta Vector) Matrix {
	m[4], m[5] = m[4]+delta.X, m[5]+delta.Y
	return m
}

// ScaledXY scales everything around a given point by the scale factor in each axis respectively.
func (m Matrix) ScaledXY(around Vector, scale Vector) Matrix {
	m[4], m[5] = m[4]-around.X, m[5]-around.Y
	m[0], m[2], m[4] = m[0]*scale.X, m[2]*scale.X, m[4]*scale.X
	m[1], m[3], m[5] = m[1]*scale.Y, m[3]*scale.Y, m[5]*scale.Y
	m[4], m[5] = m[4]+around.X, m[5]+around.Y
	return m
}

// Scaled scales everything around a given point by the scale factor.
func (m Matrix) Scaled(around Vector, scale float64) Matrix {
	return m.ScaledXY(around, V(scale, scale))
}

// Rotated rotates everything around a given point by the given angle in radians.
func (m Matrix) Rotated(around Vector, angle float64) Matrix {
	sint, cost := math.Sincos(angle)
	m[4], m[5] = m[4]-around.X, m[5]-around.Y
	m = m.Chained(Matrix{cost, sint, -sint, cost, 0, 0})
	m[4], m[5] = m[4]+around.X, m[5]+around.Y
	return m
}

// Chained adds another Matrix to this one. All tranformations by the next Matrix will be applied
// after the transformations of this Matrix.
func (m Matrix) Chained(next Matrix) Matrix {
	return Matrix{
		next[0]*m[0] + next[2]*m[1],
		next[1]*m[0] + next[3]*m[1],
		next[0]*m[2] + next[2]*m[3],
		next[1]*m[2] + next[3]*m[3],
		next[0]*m[4] + next[2]*m[5] + next[4],
		next[1]*m[4] + next[3]*m[5] + next[5],
	}
}

// Project applies all transformations added to the Matrix to a vector u and returns the result.
//
// Time complexity is O(1).
func (m Matrix) Project(v *Vector) *Vector {
	return &Vector{m[0]*nv.X + m[2]*nv.Y + m[4], m[1]*nv.X + m[3]*nv.Y + m[5]}
}

// Unproject does the inverse operation to Project.
//
// It turns out that multiplying a vector by the inverse matrix of m can be nearly-accomplished by
// subtracting the translate part of the matrix and multplying by the inverse of the top-left 2x2
// matrix, and the inverse of a 2x2 matrix is simple enough to just be inlined in the computation.
//
// Time complexity is O(1).
func (m Matrix) Unproject(v *Vector) *Vector {
	d := (m[0] * m[3]) - (m[1] * m[2])
	nv.X, nv.Y = (nv.X-m[4])/d, (nv.Y-m[5])/d
	return &Vector{nv.X*m[3] - nv.Y*m[1], nv.Y*m[0] - nv.X*m[2]}
}
*/
