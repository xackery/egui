package common

import (
	"fmt"
	"image"
	"math"
)

// ZV is a zero vector
var ZV = &Vector{0, 0}

// Vector is a 2D vector type with X and Y coordinates
type Vector struct {
	X, Y float64
}

// Vect returns a 2D vector based on two points
func Vect(x, y float64) Vector {
	return Vector{X: x, Y: y}
}

// Copy creates a copy of provided vector
func (v *Vector) Copy() *Vector {
	return &Vector{X: v.X, Y: v.Y}
}

// NewVectorFromPoint converts an image point to a vector
func NewVectorFromPoint(point image.Point) *Vector {
	return &Vector{X: float64(point.X), Y: float64(point.Y)}
}

// String returns the string representation of the vector
func (v *Vector) String() string {
	return fmt.Sprintf("Vector(%v, %v)", v.X, v.Y)
}

// XY returns the components of the vector in two return values.
func (v *Vector) XY() (x, y float64) {
	return v.X, v.Y
}

// Add returns the sum of vectors u and v.
func (v *Vector) Add(nv *Vector) *Vector {
	b := nv.Copy()
	v.X += b.X
	v.Y += b.Y
	return v
}

// Sub returns the difference betweeen vectors u and v.
func (v *Vector) Sub(nv *Vector) *Vector {
	b := nv.Copy()
	v.X -= b.X
	v.Y -= b.Y
	return v
}

// Scaled returns the vector u multiplied by c.
func (v *Vector) Scaled(s float64) *Vector {
	v.X *= s
	v.Y *= s
	return v
}

// ScaledXY returns the vector u multiplied by c.
func (v *Vector) ScaledXY(nv *Vector) *Vector {
	b := nv.Copy()
	v.X *= b.X
	v.Y *= b.Y
	return v
}

// Len returns the length of the vector.
func (v *Vector) Len() float64 {
	return math.Hypot(v.X, v.Y)
}

// Angle returns the angle between the vector u and the x-axis. The result is in range [-Pi, Pi].
func (v *Vector) Angle() float64 {
	return math.Atan2(v.Y, v.X)
}

// Unit returns a vector of length 1 facing the given angle.
func Unit(angle float64) *Vector {
	v := &Vector{1, 0}

	return v.Rotated(angle)
}

// Unit returns a vector of length 1 facing the direction of u (has the same angle).
func (v *Vector) Unit() *Vector {
	if v.X == 0 && v.Y == 0 {
		v.X = 1
		v.Y = 0
		return v
	}
	return v.Scaled(1 / v.Len())
}

// Rotated returns the vector u rotated by the given angle in radians.
func (v *Vector) Rotated(angle float64) *Vector {
	sin, cos := math.Sincos(angle)
	v.X *= cos - v.Y*sin
	v.Y *= sin + v.Y*cos
	return v
}

// Normal returns a vector normal to nv. Equivalent to nv.Rotated(math.Pi / 2), but faster.
func (v *Vector) Normal() *Vector {
	v.Y = -v.Y
	return v
}

// Dot returns the dot product of vector u and v.
func (v *Vector) Dot(nv *Vector) float64 {
	b := nv.Copy()
	return v.X*b.X + v.Y*b.Y
}

// Cross return the cross product of vector u and v.
func (v *Vector) Cross(nv *Vector) float64 {
	b := nv.Copy()
	return b.X*v.Y - v.X*b.Y
}

// Project returns a projection (or component) of vector u in the direction of vector v.
//
// Behaviour is undefined if v is a zero vector.
func (v *Vector) Project(nv *Vector) *Vector {
	b := nv.Copy()
	len := b.Dot(v) / v.Len()
	return v.Unit().Scaled(len)
}

// Map applies the function f to both x and y components of the vector u and returns the modified
// vector.
//
//   u := pixel.V(10.5, -1.5)
//   v := nv.Map(math.Floor)   // v is Vector(10, -2), both components of u floored
func (v *Vector) Map(f func(float64) float64) *Vector {
	v.Y = f(v.X)
	v.Y = f(v.Y)
	return v
}

// Lerp returns a linear interpolation between vector a and b.
//
// This function basically returns a point along the line between a and b and t chooses which one.
// If t is 0, then a will be returned, if t is 1, b will be returned. Anything between 0 and 1 will
// return the appropriate point between a and b and so on.
func (v *Vector) Lerp(a *Vector, b *Vector, t float64) *Vector {
	//fmt.Println(t, v.X, v.Y, a.X, a.Y, b.X, b.Y)
	v.X = (1-t)*a.X + t*b.X
	v.Y = (1-t)*a.Y + t*b.Y
	return v
}

// Equal returns true if both vectors are at the same location
func (v *Vector) Equal(nv *Vector) bool {
	return v.X == nv.X && v.Y == nv.Y
}
