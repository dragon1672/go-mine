package linerender

import (
	"context"
	"fmt"
	"runtime"
	"sync"
	"time"

	"github.com/dragon1672/go-mine/minecraft/utils/tickers"
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/golang/glog"
)

func init() {
	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()
}

type LineRenderer struct {
	window *glfw.Window
}

type RGB struct {
	R, G, B float32
}

func (l *LineRenderer) SetColor(red, blue, green float32) {
	gl.Color3f(red, green, blue)
}

func (l *LineRenderer) SetBackgroundColor(red, green, blue float32) {
	gl.ClearColor(red, green, blue, 0.0)
}

func (l *LineRenderer) GetSize() (width, height float32) {
	w, h := l.window.GetSize()
	return float32(w), float32(h)
}

func (l *LineRenderer) DrawLine(x1, y1, x2, y2 int) {
	gl.Begin(gl.LINES)
	gl.Vertex2f(float32(x1), float32(y1))
	gl.Vertex2f(float32(x2), float32(y2))
	gl.End()
}

func (l *LineRenderer) DrawPoint(x, y int, size float32) {
	gl.PointSize(size)
	gl.Begin(gl.POINTS)
	gl.Vertex3f(float32(x), float32(y), 0.0)
	gl.End()
}

type UpdateFunc func(t time.Time, dt time.Duration) error
type DrawFunc func(renderer *LineRenderer)

type Window struct {
	renderer *LineRenderer
	window   *glfw.Window
	update   UpdateFunc
	draw     DrawFunc
}

// Start will start calling the update and draw functions.
// This will "hang" until the window is closed.
func (w *Window) Start(ctx context.Context) error {

	glog.InfoContext(ctx, "Starting update ticker")
	updateCleanup := tickers.StartTicker(ctx, 10*time.Millisecond, func(t time.Time, dt time.Duration) (bool, error) {
		if err := w.update(t, dt); err != nil {
			return false, fmt.Errorf("update error: %v", err)
		}
		return true, nil
	})
	defer updateCleanup()

	// setup default color
	w.renderer.SetBackgroundColor(.5, .5, .5)

	glog.InfoContext(ctx, "Starting main render loop")
	for !w.window.ShouldClose() {
		select {
		case <-ctx.Done():
			return nil
		default:
			gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
			//width, height := w.renderer.GetSize()
			gl.Viewport(0, 0, 100, 100)
			gl.Ortho(0, 100, 0, 100, -1.0, 1.0)
			//gl.Ortho(-200.0, 200.0, -200.0, 200.0, -200.0, 200.0)
			w.draw(w.renderer)
			w.window.SwapBuffers()
			glfw.PollEvents()
		}
	}
	return nil
}

var getWindow = sync.OnceValues[*glfw.Window, error](func() (*glfw.Window, error) {
	if err := glfw.Init(); err != nil {
		return nil, fmt.Errorf("failed to initialize glfw: %v", err)
	}

	glfw.WindowHint(glfw.Resizable, glfw.True)
	glfw.WindowHint(glfw.ContextVersionMajor, 2)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)

	// there must be a window that is current context before gl.Init()
	window, err := glfw.CreateWindow(800, 600, "Cube", nil, nil)
	if err != nil {
		return nil, err
	}
	window.MakeContextCurrent()

	if err := gl.Init(); err != nil {
		return nil, fmt.Errorf("unable to gl.Init(): %v", err)
	}
	return window, nil
})

func MakeWindow(width, height int, update UpdateFunc, draw DrawFunc) (*Window, error) {
	w, err := getWindow()
	if err != nil {
		return nil, fmt.Errorf("error making window: %v", err)
	}
	w.SetSize(width, height)
	return &Window{
		renderer: &LineRenderer{window: w},
		window:   w,
		update:   update,
		draw:     draw,
	}, nil
}
