package renderer

import (
	"fmt"
	"runtime"
	"sync"
	"time"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

func init() {
	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()
}

const width, height = 800, 600

type Drawable interface {
	Draw(t time.Time, dt time.Duration) error
	Cleanup()
}

type Window struct {
	window     *glfw.Window
	items      []Drawable
	lastRender time.Time
}

func (w *Window) Cleanup() {
	glfw.Terminate()
	for _, r := range w.items {
		r.Cleanup()
	}
}

var getWindow = sync.OnceValues[*Window, error](func() (*Window, error) {
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
	return &Window{
		window: window,
	}, nil
})

// GetWindow creates a renderer
func GetWindow() (*Window, error) {
	return getWindow()
}

func (w *Window) drawAll(t time.Time) error {
	// All Open GL calls need to be on the main thread :(
	// Might be able to figure out a dispatcher or something to make this more sane to work with
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	dt := t.Sub(w.lastRender)

	for _, r := range w.items {
		if err := r.Draw(t, dt); err != nil {
			return err
		}
	}
	w.lastRender = t
	return nil
}

func (w *Window) RenderLoop(t time.Time) (bool, error) {
	if w.window.ShouldClose() {
		return false, nil
	}
	if err := w.drawAll(t); err != nil {
		return false, err
	}
	w.window.SwapBuffers()
	glfw.PollEvents()
	return true, nil
}

func (w *Window) AddItem(obj Drawable) {
	w.items = append(w.items, obj) // TODO thread safety
}

func (w *Window) GetWindow() *glfw.Window {
	return w.window
}
