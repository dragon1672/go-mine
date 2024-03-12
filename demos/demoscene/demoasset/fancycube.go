package demoasset

import (
	"context"
	"github.com/dragon162/go-mine/minecraft/utils/vec"
	"math/rand"
	"slices"
	"time"

	"github.com/dragon162/go-mine/minecraft/renderer/textures"
	"github.com/dragon162/go-mine/minecraft/utils/tickers"
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/golang/glog"
)

type FancyDemoCube struct {
	rotation     *vec.LinearDynoVec
	texture      uint32
	cleanupFuncs []func()
}

func (d *FancyDemoCube) Draw(t time.Time, dt time.Duration) error {
	gl.MatrixMode(gl.MODELVIEW)
	gl.LoadIdentity()
	gl.Translatef(0, 0, -3.0)
	rotation := d.rotation.Get(t)
	gl.Rotatef(float32(rotation.X), 1, 0, 0)
	gl.Rotatef(float32(rotation.Y), 0, 1, 0)

	gl.BindTexture(gl.TEXTURE_2D, d.texture)

	gl.Color4f(1, 1, 1, 1)

	gl.Begin(gl.QUADS)

	gl.Normal3f(0, 0, 1)
	gl.TexCoord2f(0, 0)
	gl.Vertex3f(-1, -1, 1)
	gl.TexCoord2f(1, 0)
	gl.Vertex3f(1, -1, 1)
	gl.TexCoord2f(1, 1)
	gl.Vertex3f(1, 1, 1)
	gl.TexCoord2f(0, 1)
	gl.Vertex3f(-1, 1, 1)

	gl.Normal3f(0, 0, -1)
	gl.TexCoord2f(1, 0)
	gl.Vertex3f(-1, -1, -1)
	gl.TexCoord2f(1, 1)
	gl.Vertex3f(-1, 1, -1)
	gl.TexCoord2f(0, 1)
	gl.Vertex3f(1, 1, -1)
	gl.TexCoord2f(0, 0)
	gl.Vertex3f(1, -1, -1)

	gl.Normal3f(0, 1, 0)
	gl.TexCoord2f(0, 1)
	gl.Vertex3f(-1, 1, -1)
	gl.TexCoord2f(0, 0)
	gl.Vertex3f(-1, 1, 1)
	gl.TexCoord2f(1, 0)
	gl.Vertex3f(1, 1, 1)
	gl.TexCoord2f(1, 1)
	gl.Vertex3f(1, 1, -1)

	gl.Normal3f(0, -1, 0)
	gl.TexCoord2f(1, 1)
	gl.Vertex3f(-1, -1, -1)
	gl.TexCoord2f(0, 1)
	gl.Vertex3f(1, -1, -1)
	gl.TexCoord2f(0, 0)
	gl.Vertex3f(1, -1, 1)
	gl.TexCoord2f(1, 0)
	gl.Vertex3f(-1, -1, 1)

	gl.Normal3f(1, 0, 0)
	gl.TexCoord2f(1, 0)
	gl.Vertex3f(1, -1, -1)
	gl.TexCoord2f(1, 1)
	gl.Vertex3f(1, 1, -1)
	gl.TexCoord2f(0, 1)
	gl.Vertex3f(1, 1, 1)
	gl.TexCoord2f(0, 0)
	gl.Vertex3f(1, -1, 1)

	gl.Normal3f(-1, 0, 0)
	gl.TexCoord2f(0, 0)
	gl.Vertex3f(-1, -1, -1)
	gl.TexCoord2f(1, 0)
	gl.Vertex3f(-1, -1, 1)
	gl.TexCoord2f(1, 1)
	gl.Vertex3f(-1, 1, 1)
	gl.TexCoord2f(0, 1)
	gl.Vertex3f(-1, 1, -1)

	gl.End()
	return nil
}

func (d *FancyDemoCube) tick(t time.Time) {
	r := func() float64 { return 200.0 * (rand.Float64()*2.0 - 1.0) }
	newRot := vec.Vec3{
		X: r(),
		Y: r(),
	}
	d.rotation.Set(t, d.rotation.Get(t), newRot)
}

func (d *FancyDemoCube) StartTicks(ctx context.Context) {
	spinCleanup := tickers.StartTicker(ctx, 3*time.Second, func(t time.Time, dt time.Duration) (bool, error) {
		d.tick(t)
		return true, nil
	})
	d.cleanupFuncs = append(d.cleanupFuncs, spinCleanup)
}

func (d *FancyDemoCube) Cleanup() {
	gl.DeleteTextures(1, &d.texture)
	// Grap a copy of the cleanup functions and clear the list to avoid duplicate calls
	// Note this is not thread safe
	cleanups := slices.Clone(d.cleanupFuncs)
	d.cleanupFuncs = nil
	for _, cleanup := range cleanups {
		cleanup()
	}
}

func MakeFancyCube() (*FancyDemoCube, error) {
	glog.Info("Creating fancy demo cube with hard coded texture")
	texture, err := textures.LoadTextureFile("square.png")
	if err != nil {
		return nil, err
	}
	return &FancyDemoCube{
		texture:  textures.LoadTextureToGPU(texture),
		rotation: vec.MakeLinearDynoVec(time.Now(), vec.Vec3{}, vec.Vec3{}),
	}, nil
}
