package blocks

type Block interface {
	DirtyGen() bool
}

type SimpleBlockType int64

const (
	Air SimpleBlockType = iota
	Grass
	Sand
	Dirt
	Stone
	Leaves
	Wood
	Flower
)

func (s SimpleBlockType) String() string {
	switch s {
	case Air:
		return "Air"
	case Grass:
		return "Grass"
	case Sand:
		return "Sand"
	case Dirt:
		return "Dirt"
	case Stone:
		return "Stone"
	case Leaves:
		return "Leaves"
	case Wood:
		return "Wood"
	case Flower:
		return "Flower"
	}
	return "unknown"
}

type SimpleBlock struct {
	Dirtied bool
	T       SimpleBlockType
}

func (s *SimpleBlock) DirtyGen() bool {
	return s.Dirtied
}
