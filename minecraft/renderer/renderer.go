package renderer

import (
	"fmt"
	"runtime"
	"time"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

func init() {
	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()
}

const width, height = 800, 600

// TODO: yush
var renderer *Renderer // Highlander there can only be one!

type Renderable interface {
	Draw(t time.Time) error
	Cleanup()
}

type Renderer struct {
	window      *glfw.Window
	renderables []Renderable
}

func (r *Renderer) Cleanup() {
	glfw.Terminate()
	for _, r := range r.renderables {
		r.Cleanup()
	}
}

// InitRenderer make da things
// TODO: only call once
func InitRenderer() (*Renderer, error) {
	if renderer == nil {
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
		renderer = &Renderer{
			window: window,
		}
	}
	return renderer, nil
}

func (r *Renderer) drawAll(t time.Time) error {
	// All Open GL calls need to be on the main thread :(
	// Might be able to figure out a dispatcher or something to make this more sane to work with
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	for _, r := range r.renderables {
		if err := r.Draw(t); err != nil {
			return err
		}
	}
	return nil
}

func (r *Renderer) RenderLoop(t time.Time) (bool, error) {
	if r.window.ShouldClose() {
		return false, nil
	}
	if err := r.drawAll(t); err != nil {
		return false, err
	}
	r.window.SwapBuffers()
	glfw.PollEvents()
	return true, nil
}

func (r *Renderer) AddRenderable(obj Renderable) {
	r.renderables = append(r.renderables, obj) // TODO thread safety
}