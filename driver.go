package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/dragon1672/go-mine/demos/demoscene"
	"github.com/golang/glog"
)

func main() {
	flag.Parse()
	ctx := context.Background()
	fmt.Println("Hello World!")
	//glhfdemo.DemoMain()
	demoscene.BadMain(ctx)
	glog.Flush()
}
