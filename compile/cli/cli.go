package main

import (
	"flag"
	"io/ioutil"
	"os"

	"github.com/soloworks/go-netlinx/apw"
	"github.com/soloworks/go-netlinx/compile"
)

type myargs struct {
	Source string
	Dest   string
	Root   string
}

var args myargs

func main() {
	// Get Command Line Variables
	flag.StringVar(&args.Source, "Source", "", "Source APW File")
	flag.StringVar(&args.Dest, "Dest", "compile.cfg", "Destination CFG File")
	flag.StringVar(&args.Root, "Root", ".", "Root Directory")
	flag.Parse()

	// Load in the core APW file
	a, err := apw.LoadAPW(args.Source)
	if err != nil {
		println(`Error Loading APW File: "` + args.Source + `"`)
		println(err)
		os.Exit(1)
	}

	// Process and generate the .cfg
	b := compile.GenerateCFG(*a, args.Root, "", true)

	// Output to File
	err = ioutil.WriteFile(args.Dest, b, 0644)
	if err != nil {
		println(`Error Writing CFG File: "` + args.Dest + `"`)
		println(err)
		os.Exit(1)
	}
}
