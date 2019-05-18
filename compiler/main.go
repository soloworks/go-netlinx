package main

import (
	"flag"
	"github.com/soloworks/go-netlinx/apw"
)

type myargs struct {
	Source string
	Dest   string
}

var args myargs

func main() {
	// Get Command Line Variables 
	flag.StringVar(&args.Source, "Source", "", "Source APW File")
	flag.StringVar(&args.Dest, "Dest", "", "Destination CFG File")
	flag.Parse()

	// Load in the core APW file
	apw, err := apw.LoadAPWInfo(args.Source) x

	if err != nil {
		println("Error Loading CoreAPW:" + err.Error()) 
	} 
}
