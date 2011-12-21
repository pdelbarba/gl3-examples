package main

import (
	"math"
)

// ### Vector3 ###

type Vector3 []float32

func NewVector3(x, y, z float32) *Vector3 {
	return &Vector3{x, y, z}
}

func (v Vector3) LengthSqrt() float32 {
	return v[0]*v[0] + v[1]*v[1] + v[2]*v[2]
}

// Returns the length of the vector.
func (v *Vector3) Magnitude() float32 {
	return float32(math.Sqrt(float64(v.LengthSqrt())))
}

func (v Vector3) Normalize() {
	l := 1.0 / v.Magnitude()
	v[0] *= l
	v[1] *= l
	v[2] *= l
}

func (v Vector3) Normalized() *Vector3 {
	l := 1.0 / v.Magnitude()
	return &Vector3{v[0] * l, v[1] * l, v[2] * l}
}
