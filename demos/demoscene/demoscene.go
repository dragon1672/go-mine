package demoscene

import (
	"context"
	"github.com/dragon162/go-mine/demos/demoscene/demoasset"
	"time"

	"github.com/dragon162/go-mine/minecraft/renderer"
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/golang/glog"
)

func setupScene(r *renderer.Window) {
	glog.Info("Setting up general area with some lights")
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
	width, height := r.GetWindow().GetSize()
	f := float64(width)/float64(height) - 1.0
	gl.Frustum(-1-f, 1+f, -1, 1, 1.0, 10.0)
	gl.MatrixMode(gl.MODELVIEW)
	gl.LoadIdentity()
}

func BadGameLoop(w *renderer.Window) error {
	glog.Info("Starting main render loop")
	for {
		ok, err := w.RenderLoop(time.Now())
		if err != nil {
			return err
		}
		if !ok {
			break
		}
	}
	return nil
}

func BadMain(ctx context.Context) {
	w, err := renderer.GetWindow()
	if err != nil {
		glog.Fatalf("error creating game window: %v", err)
	}
	defer w.Cleanup()

	setupScene(w)

	//cube, err := demoasset.MakeCube()
	cube, err := demoasset.MakeFancyCube()
	if err != nil {
		glog.Fatalf("error making cube: %v", err)
	}

	glog.Info("Add cube to window to be rendered")
	w.AddItem(cube)

	glog.Info("Start 'gameloop' of cube so it will update")
	cube.StartTicks(ctx)

	if err := BadGameLoop(w); err != nil {
		glog.Fatalf("Game Loop Err: %v", err)
	}
}
