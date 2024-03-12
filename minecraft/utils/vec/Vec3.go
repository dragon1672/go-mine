package vec

type Vec3 struct {
	X, Y, Z float64
}

func (v Vec3) Add(t Vec3) Vec3 {
	return Vec3{X: v.X + t.X, Y: v.Y + t.Y, Z: v.Z + t.Z}
}

func (v Vec3) Mul(f float64) Vec3 {
	return Vec3{X: v.X * f, Y: v.Y * f, Z: v.Z * f}
}
