package demoasset

import (
	"context"
	"math/rand"
	"slices"
	"time"

	"github.com/dragon162/go-mine/minecraft/renderer/textures"
	"github.com/dragon162/go-mine/minecraft/utils/tickers"
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/golang/glog"
)

type DemoCube struct {
	rotationX, rotationY float64
	texture              uint32
	xSpeed, ySpeed       float64
	cleanupFuncs         []func()
}

func (d *DemoCube) Draw(t time.Time, dt time.Duration) error {
	gl.MatrixMode(gl.MODELVIEW)
	gl.LoadIdentity()
	gl.Translatef(0, 0, -3.0)
	gl.Rotatef(float32(d.rotationX), 1, 0, 0)
	gl.Rotatef(float32(d.rotationY), 0, 1, 0)

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

func (d *DemoCube) tick() {
	d.xSpeed = 200.0 * (rand.Float64()*2.0 - 1.0)
	d.ySpeed = 200.0 * (rand.Float64()*2.0 - 1.0)
}

func (d *DemoCube) updateTick(t time.Time, dt time.Duration) {
	d.rotationX += d.xSpeed * dt.Seconds()
	d.rotationY += d.xSpeed * dt.Seconds()
}

func (d *DemoCube) StartTicks(ctx context.Context) {
	spinCleanup := tickers.StartTicker(ctx, 1*time.Second, func(t time.Time, dt time.Duration) (bool, error) {
		d.tick()
		return true, nil
	})
	d.cleanupFuncs = append(d.cleanupFuncs, spinCleanup)

	updateCleanup := tickers.StartTicker(ctx, 1*time.Millisecond, func(t time.Time, dt time.Duration) (bool, error) {
		d.updateTick(t, dt)
		return true, nil
	})

	d.cleanupFuncs = append(d.cleanupFuncs, updateCleanup)
}

func (d *DemoCube) Cleanup() {
	gl.DeleteTextures(1, &d.texture)
	// Grap a copy of the cleanup functions and clear the list to avoid duplicate calls
	// Note this is not thread safe
	cleanups := slices.Clone(d.cleanupFuncs)
	d.cleanupFuncs = nil
	for _, cleanup := range cleanups {
		cleanup()
	}
}

func MakeCube() (*DemoCube, error) {
	glog.Info("Creating demo cube with hard coded texture")
	texture, err := textures.LoadTextureFile("square.png")
	if err != nil {
		return nil, err
	}
	return &DemoCube{
		texture:   textures.LoadTextureToGPU(texture),
		rotationX: 0,
		rotationY: 0,
	}, nil
}
