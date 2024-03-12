package main

import (
	"flag"
	"fmt"
	"github.com/dragon162/go-mine/demos/demoscene"
	"github.com/golang/glog"
)

func main() {
	flag.Parse()
	fmt.Println("Hello World!")
	//glhfdemo.DemoMain()
	demoscene.BadMain()
	glog.Flush()
}
