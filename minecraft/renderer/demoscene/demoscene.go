package demoscene

import (
	"github.com/dragon162/go-mine/minecraft/renderer"
	"github.com/dragon162/go-mine/minecraft/renderer/demoasset"
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/golang/glog"
	"time"
)

const width, height = 800, 600

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
	f := float64(width)/height - 1 // TODO change ratio to be off live dimensions
	gl.Frustum(-1-f, 1+f, -1, 1, 1.0, 10.0)
	gl.MatrixMode(gl.MODELVIEW)
	gl.LoadIdentity()
}

func BadGameLoop(r *renderer.Renderer) error {
	var lastTime time.Time
	for {
		t := time.Now()
		dt := t.Sub(lastTime)
		lastTime = t
		ok, err := r.RenderLoop(t, dt)
		if err != nil {
			return err
		}
		if !ok {
			break
		}
		glog.Infof("Yay! @ %v", t)
	}
	return nil
}

func BadMain() {
	rendy, err := renderer.GetRenderer()
	if err != nil {
		glog.Fatalf("error creating renderer: %v", err)
	}
	defer rendy.Cleanup()

	setupScene()

	cube, err := demoasset.MakeCube()
	if err != nil {
		glog.Fatalf("error making cube: %v", err)
	}

	rendy.AddRenderable(cube)

	cube.StartTicks()

	if err := BadGameLoop(rendy); err != nil {
		glog.Fatalf("Game Loop Err: %v", err)
	}
}
