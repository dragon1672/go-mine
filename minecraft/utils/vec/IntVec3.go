package vec

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
