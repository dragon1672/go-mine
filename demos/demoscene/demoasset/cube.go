package demoasset

import (
	"fmt"
	"image"
	"image/draw"
	_ "image/png"
	"math/rand"
	"os"
	"slices"
	"time"

	"github.com/dragon162/go-mine/minecraft/utils/tickers"
	"github.com/go-gl/gl/v2.1/gl"
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

func (d *DemoCube) StartTicks() {
	spinCleanup := tickers.StartTicker(1*time.Second, func(t time.Time, dt time.Duration) (bool, error) {
		d.tick()
		return true, nil
	})
	d.cleanupFuncs = append(d.cleanupFuncs, spinCleanup)

	updateCleanup := tickers.StartTicker(1*time.Millisecond, func(t time.Time, dt time.Duration) (bool, error) {
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
	texture, err := loadTextureFile("square.png")
	if err != nil {
		return nil, err
	}
	return &DemoCube{
		texture:   loadTextureToGPU(texture),
		rotationX: 0,
		rotationY: 0,
	}, nil
}

func loadTextureFile(file string) (*image.RGBA, error) {
	imgFile, err := os.Open(file)
	if err != nil {
		return nil, fmt.Errorf("texture %q not found on disk: %v\n", file, err)
	}
	img, _, err := image.Decode(imgFile)
	if err != nil {
		return nil, fmt.Errorf("error decoding image: %v", err)
	}

	rgba := image.NewRGBA(img.Bounds())
	if rgba.Stride != rgba.Rect.Size().X*4 {
		return nil, fmt.Errorf("unsupported stride %d, expected %d", rgba.Stride, rgba.Rect.Size().X*4)
	}
	draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)
	return rgba, nil
}

func loadTextureToGPU(rgba *image.RGBA) uint32 {
	var texture uint32
	gl.Enable(gl.TEXTURE_2D)
	gl.GenTextures(1, &texture)
	gl.BindTexture(gl.TEXTURE_2D, texture)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexImage2D(
		gl.TEXTURE_2D,
		0,
		gl.RGBA,
		int32(rgba.Rect.Size().X),
		int32(rgba.Rect.Size().Y),
		0,
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		gl.Ptr(rgba.Pix))

	return texture
}
