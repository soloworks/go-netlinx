package main

import (
	"flag"
	"io/ioutil"
	"os"

	"github.com/soloworks/go-netlinx/compilelog"
)

type myargs struct {
	Source string
	Dest   string
	Root   string
}

var args myargs

func main() {
	// Get Command Line Variables
	flag.StringVar(&args.Source, "Source", "", "Source Log File")
	flag.StringVar(&args.Root, "Root", "", "Root Directory for log")
	flag.StringVar(&args.Dest, "Dest", "clean.log", "Destination Log File")
	flag.Parse()

	// Load in the core APW file
	b, err := ioutil.ReadFile(args.Source)
	if err != nil {
		println(`Error Loading Log File: "` + args.Source + `"`)
		println(err.Error())
		os.Exit(1)
	}

	// Process the log
	result, err := compilelog.Process(b, args.Root)

	// Output to File
	err = ioutil.WriteFile(args.Dest, result, 0644)
	if err != nil {
		println(`Error Writing Log File: "` + args.Dest + `"`)
		println(err)
		os.Exit(1)
	}

}
