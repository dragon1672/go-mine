package worldgen

import (
	"fmt"
	"github.com/dragon162/go-mine/minecraft/world"
	"github.com/dragon162/go-mine/minecraft/world/blocks"
	"testing"
)

func TestBaseGen(t *testing.T) {
	g := New(42)
	size := 5
	var layers [][][]blocks.SimpleBlockType
	for y := 0; y < 20; y++ {
		var layer [][]blocks.SimpleBlockType
		for x := -size; x < size; x++ {
			var slice []blocks.SimpleBlockType
			for z := -size; z < size; z++ {
				p := world.Vec3Of(x, y, z)
				slice = append(slice, g.baseGen(p))
			}
			layer = append(layer, slice)
		}
		layers = append(layers, layer)
	}
	fmt.Print(layers)
}
