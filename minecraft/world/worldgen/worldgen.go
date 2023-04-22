package worldgen

import (
	"github.com/dragon162/go-mine/minecraft/world"
	"github.com/dragon162/go-mine/minecraft/world/blocks"
	"github.com/ojrac/opensimplex-go"
)

type WorldGenny struct {
	sim opensimplex.Noise
}

func (w *WorldGenny) noise2(x, y float32, octaves int, persistence, lacunarity float32) float32 {
	var (
		freq  float32 = 1
		amp   float32 = 1
		max   float32 = 1
		total         = w.sim.Eval2(float64(x), float64(y))
	)
	for i := 0; i < octaves; i++ {
		freq *= lacunarity
		amp *= persistence
		max += amp
		total += w.sim.Eval2(float64(x*freq), float64(y*freq)) * float64(amp)
	}
	return (1 + float32(total)/max) / 2
}

func (w *WorldGenny) noise3(x, y, z float32, octaves int, persistence, lacunarity float32) float32 {
	var (
		freq  float32 = 1
		amp   float32 = 1
		max   float32 = 1
		total         = w.sim.Eval3(float64(x), float64(y), float64(z))
	)
	for i := 0; i < octaves; i++ {
		freq *= lacunarity
		amp *= persistence
		max += amp
		total += w.sim.Eval3(float64(x*freq), float64(y*freq), float64(z*freq)) * float64(amp)
	}
	return (1 + float32(total)/max) / 2
}

// https://github.com/icexin/gocraft/blob/50535c9c92c7156b161e0a06ea3eedf12e11a37b/chunk.go#L15
func (w *WorldGenny) rawGen(pos world.Vec3) blocks.SimpleBlockType {
	f := w.noise2(float32(pos.X)*0.01, float32(pos.Z)*0.01, 4, 0.5, 2)
	g := w.noise2(float32(-pos.X)*0.01, float32(-pos.Z)*0.01, 2, 0.9, 2)
	groundLevel := int(f * float32(int(g*32+16))) // Top of ground
	if groundLevel <= 12 {
		groundLevel = 12 // baseline to prevent too much depth
	}

	if pos.Y > groundLevel {
		return blocks.Air // clear the skys
	}

	if pos.Y < groundLevel-5 {
		return blocks.Stone // TODO make the underground interesting
	}

	// grass?
	if beGrass := w.noise2(-float32(pos.X)*0.1, float32(pos.Z)*0.1, 4, 0.8, 2) > 0.6; beGrass {
		if pos.Y == groundLevel {
			return blocks.Grass
		}
		// no else since sky already covered
		return blocks.Dirt
	} else {
		// no checks since sky already covered
		return blocks.Sand
	}
}

func (w *WorldGenny) baseGen(pos world.Vec3) blocks.SimpleBlockType {
	if w.rawGen(pos.Down()) == blocks.Grass {
		// flowers
		if w.noise2(float32(pos.X)*0.05, float32(-pos.Z)*0.05, 4, 0.8, 2) > 0.7 {
			return blocks.Flower
		}
	}
	return w.rawGen(pos) // default to nothing interesting
}

func (w *WorldGenny) GetGen(pos world.Vec3) blocks.Block {
	return &blocks.SimpleBlock{
		Dirtied: false,
		T:       w.baseGen(pos),
	}
}

func New(seed int64) *WorldGenny {
	return &WorldGenny{
		sim: opensimplex.New(seed),
	}
}
