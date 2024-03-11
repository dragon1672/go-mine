package world

type IntVec3 struct {
	X, Y, Z int
}

func (v IntVec3) Left() IntVec3 {
	return IntVec3{v.X - 1, v.Y, v.Z}
}
func (v IntVec3) Right() IntVec3 {
	return IntVec3{v.X + 1, v.Y, v.Z}
}
func (v IntVec3) Up() IntVec3 {
	return IntVec3{v.X, v.Y + 1, v.Z}
}
func (v IntVec3) Down() IntVec3 {
	return IntVec3{v.X, v.Y - 1, v.Z}
}
func (v IntVec3) Front() IntVec3 {
	return IntVec3{v.X, v.Y, v.Z + 1}
}
func (v IntVec3) Back() IntVec3 {
	return IntVec3{v.X, v.Y, v.Z - 1}
}

type Vec3 struct {
	X, Y, Z float64
}

func (v Vec3) Add(t Vec3) Vec3 {
	return Vec3{X: v.X + t.X, Y: v.Y + t.Y, Z: v.Z + t.Z}
}

func (v Vec3) Mul(f float64) Vec3 {
	return Vec3{X: v.X * f, Y: v.Y * f, Z: v.Z * f}
}
