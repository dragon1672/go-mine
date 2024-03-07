package renderer

import (
	"fmt"
	"image"
	"image/draw"
	_ "image/png"
	"log"
	"math/rand"
	"os"
	"runtime"
	"time"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

// TODO: yush
var renderer *Renderer // Highlander there can only be one!

type Renderable interface {
	Draw() error
	Cleanup() error
	Tick() error // TODO: move off of renderable
}

type Renderer struct {
	window      *glfw.Window
	renderables []Renderable
	ticker      *time.Ticker
	tickerExit  chan struct{}
}

func (r *Renderer) Cleanup() error {
	glfw.Terminate()
	for _, r := range r.renderables {
		if err := r.Cleanup(); err != nil {
			return err
		}
	}
	r.ticker.Stop()
	r.tickerExit <- struct{}{}
	return nil
}

func (r *Renderer) tick() {
	for _, r := range r.renderables {
		if err := r.Tick(); err != nil {
			log.Fatal(err)
		}
	}
	r.window.SwapBuffers()
}

func (r *Renderer) StartTick() {
	// TODO make sure to prevent double call
	r.ticker = time.NewTicker(1 * time.Second)
	r.tickerExit = make(chan struct{})
	go func() {
		for {
			select {
			case <-r.tickerExit:
				return
			case <-r.ticker.C:
				r.tick()
			}
		}
	}()
}

// InitRenderer make da things
// TODO: only call once
func InitRenderer() (*Renderer, error) {
	if err := glfw.Init(); err != nil {
		return nil, fmt.Errorf("failed to initialize glfw: %v", err)
	}

	glfw.WindowHint(glfw.Resizable, glfw.True)
	glfw.WindowHint(glfw.ContextVersionMajor, 2)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)

	// there must be a window that is current context before gl.Init()
	window, err := glfw.CreateWindow(width, height, "Cube", nil, nil)
	if err != nil {
		return nil, err
	}
	window.MakeContextCurrent()

	if err := gl.Init(); err != nil {
		return nil, fmt.Errorf("unable to gl.Init(): %v", err)
	}
	setupScene()
	return &Renderer{
		window: window,
	}, nil
}

func setupScene() {
	gl.Enable(gl.DEPTH_TEST)
	gl.Enable(gl.LIGHTING)

	gl.ClearColor(0.5, 0.5, 0.5, 0.0)
	gl.ClearDepth(1)
	gl.DepthFunc(gl.LEQUAL)

	ambient := []float32{0.5, 0.5, 0.5, 1}
	diffuse := []float32{1, 1, 1, 1}
	lightPosition := []float32{-5, 5, 10, 0}
	gl.Lightfv(gl.LIGHT0, gl.AMBIENT, &ambient[0])
	gl.Lightfv(gl.LIGHT0, gl.DIFFUSE, &diffuse[0])
	gl.Lightfv(gl.LIGHT0, gl.POSITION, &lightPosition[0])
	gl.Enable(gl.LIGHT0)

	gl.MatrixMode(gl.PROJECTION)
	gl.LoadIdentity()
	f := float64(width)/height - 1
	gl.Frustum(-1-f, 1+f, -1, 1, 1.0, 10.0)
	gl.MatrixMode(gl.MODELVIEW)
	gl.LoadIdentity()
}

func (r *Renderer) drawAll() {
	// All Open GL calls need to be on the main thread :(
	// Might be able to figure out a dispatcher or something to make this more sane to work with
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	for _, r := range r.renderables {
		if err := r.Draw(); err != nil {
			log.Fatal(err) // TODO error handling
		}
	}
}

func (r *Renderer) BadGameLoop() {
	r.StartTick()
	for !r.window.ShouldClose() {
		r.drawAll()
		r.window.SwapBuffers()
		glfw.PollEvents()
	}
}

func (r *Renderer) AddRenderable(obj Renderable) {
	r.renderables = append(r.renderables, obj) // TODO thread safety
}

type demoCube struct {
	rotationX, rotationY float32
	texture              uint32
	xSpeed, ySpeed       float32
}

func (d *demoCube) Draw() error {
	gl.MatrixMode(gl.MODELVIEW)
	gl.LoadIdentity()
	gl.Translatef(0, 0, -3.0)
	gl.Rotatef(d.rotationX, 1, 0, 0)
	gl.Rotatef(d.rotationY, 0, 1, 0)

	// TODO: make draw immutable
	d.rotationX += d.xSpeed
	d.rotationY += d.xSpeed

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

func (d *demoCube) Tick() error {
	d.xSpeed = 2.0 * (rand.Float32()*2.0 - 1.0)
	d.ySpeed = 2.0 * (rand.Float32()*2.0 - 1.0)
	return nil
}

func (d *demoCube) Cleanup() error {
	gl.DeleteTextures(1, &d.texture)
	return nil
}

func MakeDemoCube() (Renderable, error) {
	return &demoCube{
		texture:   newTexture("square.png"),
		rotationX: 0,
		rotationY: 0,
	}, nil
}

const width, height = 800, 600

func init() {
	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()
}

func BadMain() {
	rendy, err := InitRenderer()
	if err != nil {
		log.Fatal(err)
	}
	defer rendy.Cleanup()

	cube, err := MakeDemoCube()
	if err != nil {
		log.Fatal(err)
	}
	rendy.AddRenderable(cube)

	rendy.BadGameLoop()
}

func newTexture(file string) uint32 {
	imgFile, err := os.Open(file)
	if err != nil {
		log.Fatalf("texture %q not found on disk: %v\n", file, err)
	}
	img, _, err := image.Decode(imgFile)
	if err != nil {
		panic(err)
	}

	rgba := image.NewRGBA(img.Bounds())
	if rgba.Stride != rgba.Rect.Size().X*4 {
		panic("unsupported stride")
	}
	draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)

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
