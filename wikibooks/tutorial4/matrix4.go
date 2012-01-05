package main

import (
	"math"
)

type Matrix4 []float32

func MakeIdentity() Matrix4 {
	return Matrix4{
		1, 0, 0, 0,
		0, 1, 0, 0,
		0, 0, 1, 0,
		0, 0, 0, 1}
}

func MakeScaleMatrix(x, y, z float32) Matrix4 {
	return Matrix4{
		x, 0, 0, 0,
		0, y, 0, 0,
		0, 0, z, 0,
		0, 0, 0, 1}
}

func MakeTranslationMatrix(x, y, z float32) Matrix4 {
	return Matrix4{
		1, 0, 0, x,
		0, 1, 0, y,
		0, 0, 1, z,
		0, 0, 0, 1}
}

func MakeRotationMatrix(x, y, z float32) Matrix4 {
	xM := MakeXRotationMatrix(x)
	yM := MakeYRotationMatrix(y)
	zM := MakeZRotationMatrix(z)
	return (xM.Multiply(yM)).Multiply(zM)
}

func MakeXRotationMatrix(amount float32) Matrix4 {
	cos := float32(math.Cos(float64(amount)))
	sin := float32(math.Sin(float64(amount)))
	return Matrix4{
		1, 0, 0, 0,
		0, cos, -sin, 0,
		0, sin, cos, 0,
		0, 0, 0, 1}
}

func MakeYRotationMatrix(amount float32) Matrix4 {
	cos := float32(math.Cos(float64(amount)))
	sin := float32(math.Sin(float64(amount)))
	return Matrix4{
		cos, 0, sin, 0,
		0, 1, 0, 0,
		-sin, 0, cos, 0,
		0, 0, 0, 1}
}

func MakeZRotationMatrix(amount float32) Matrix4 {
	cos := float32(math.Cos(float64(amount)))
	sin := float32(math.Sin(float64(amount)))
	return Matrix4{
		cos, -sin, 0, 0,
		sin, cos, 0, 0,
		0, 0, 1, 0,
		0, 0, 0, 1}
}

func MakePerspectiveMatrix(fovy, aspect, zNear, zFar float32) Matrix4 {
	f := 1 / float32(math.Tan(float64(fovy/2)))
	a := 1 / (zNear - zFar)
	return Matrix4{
		f / aspect, 0, 0, 0,
		0, f, 0, 0,
		0, 0, (zFar + zNear) * a, 2 * zFar * zNear * a,
		0, 0, -1, 0,
	}
}

func MakeLookAtMatrix(eye, center, up Vector3) Matrix4 {
	f := center.Sub(eye).Normalized()
	u := up.Normalized()
	s := f.Cross(u)
	u = s.Cross(f)
	t := MakeTranslationMatrix(-eye[0], -eye[1], -eye[2])
	return Matrix4{
		s[0], s[1], s[2], 0,
		u[0], u[1], u[2], 0,
		-f[0], -f[1], -f[2], 0,
		0, 0, 0, 1}.Multiply(t)
}

func (m1 Matrix4) Multiply(m2 Matrix4) Matrix4 {
	return Matrix4{
		m1[0]*m2[0] + m1[1]*m2[4] + m1[2]*m2[8] + m1[3]*m2[12],
		m1[0]*m2[1] + m1[1]*m2[5] + m1[2]*m2[9] + m1[3]*m2[13],
		m1[0]*m2[2] + m1[1]*m2[6] + m1[2]*m2[10] + m1[3]*m2[14],
		m1[0]*m2[3] + m1[1]*m2[7] + m1[2]*m2[11] + m1[3]*m2[15],

		m1[4]*m2[0] + m1[5]*m2[4] + m1[6]*m2[8] + m1[7]*m2[12],
		m1[4]*m2[1] + m1[5]*m2[5] + m1[6]*m2[9] + m1[7]*m2[13],
		m1[4]*m2[2] + m1[5]*m2[6] + m1[6]*m2[10] + m1[7]*m2[14],
		m1[4]*m2[3] + m1[5]*m2[7] + m1[6]*m2[11] + m1[7]*m2[15],

		m1[8]*m2[0] + m1[9]*m2[4] + m1[10]*m2[8] + m1[11]*m2[12],
		m1[8]*m2[1] + m1[9]*m2[5] + m1[10]*m2[9] + m1[11]*m2[13],
		m1[8]*m2[2] + m1[9]*m2[6] + m1[10]*m2[10] + m1[11]*m2[14],
		m1[8]*m2[3] + m1[9]*m2[7] + m1[10]*m2[11] + m1[11]*m2[15],

		m1[12]*m2[0] + m1[13]*m2[4] + m1[14]*m2[8] + m1[15]*m2[12],
		m1[12]*m2[1] + m1[13]*m2[5] + m1[14]*m2[9] + m1[15]*m2[13],
		m1[12]*m2[2] + m1[13]*m2[6] + m1[14]*m2[10] + m1[15]*m2[14],
		m1[12]*m2[3] + m1[13]*m2[7] + m1[14]*m2[11] + m1[15]*m2[15],
	}
}

func (m Matrix4) Transposed() Matrix4 {
	return Matrix4{
		m[0], m[4], m[8], m[12],
		m[1], m[5], m[9], m[13],
		m[2], m[6], m[10], m[14],
		m[3], m[7], m[11], m[15],
	}
}
